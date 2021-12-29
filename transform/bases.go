package transform

import (
	"fmt"
	"io"
	"log"
	"path/filepath"

	"github.com/schollz/progressbar/v3"
)

func (c *company) porte(v string) error {
	i, err := toInt(v)
	if err != nil {
		return fmt.Errorf("error trying to parse CodigoPorte %s: %w", v, err)
	}

	var s string
	switch *i {
	case 0:
		s = "NÃO INFORMADO"
	case 1:
		s = "MICRO EMPRESA"
	case 3:
		s = "EMPRESA DE PEQUENO PORTE"
	case 5:
		s = "DEMAIS"
	}

	c.CodigoPorte = i
	if s != "" {
		c.Porte = &s
	}
	return nil
}

func (c *company) base(r []string, l *lookups) error {
	c.RazaoSocial = r[1]
	codigoNaturezaJuridica, err := toInt(r[2])
	if err != nil {
		return fmt.Errorf("error trying to parse CodigoNaturezaJuridica %s: %w", r[2], err)
	}
	c.CodigoNaturezaJuridica = codigoNaturezaJuridica
	qualificacaoDoResponsavel, err := toInt(r[3])
	if err != nil {
		return fmt.Errorf("error trying to parse QualificacaoDoResponsavel %s: %w", r[3], err)
	}
	c.QualificacaoDoResponsavel = qualificacaoDoResponsavel
	capitalSocial, err := toFloat(r[4])
	if err != nil {
		return fmt.Errorf("error trying to parse CapitalSocial %s: %w", r[4], err)
	}
	c.CapitalSocial = capitalSocial
	err = c.porte(r[5])
	if err != nil {
		return fmt.Errorf("error trying to parse Porte %s: %w", r[5], err)
	}
	enteFederativoResponsavel, err := toInt(r[6])
	if err != nil {
		return fmt.Errorf("error trying to parse EnteFederativoResponsavel%s: %w", r[6], err)
	}
	c.EnteFederativoResponsavel = enteFederativoResponsavel
	natures := l.natures[*c.CodigoNaturezaJuridica]
	if natures != "" {
		c.NaturezaJuridica = &natures
	}
	return nil
}

func addBase(l *lookups, dir string, r []string) error {
	b, err := pathForBaseCNPJ(r[0])
	if err != nil {
		return fmt.Errorf("error getting the path for %s: %w", r[0], err)
	}
	ls, err := filepath.Glob(filepath.Join(dir, b, "*.json"))
	if err != nil {
		return fmt.Errorf("error in the glob pattern: %w", err)
	}
	if len(ls) == 0 {
		log.Output(2, fmt.Sprintf("No JSON file found for CNPJ base %s", r[0]))
		return nil
	}
	for _, f := range ls {
		c, err := companyFromJSON(f)
		if err != nil {
			return fmt.Errorf("error reading company from %s: %w", f, err)
		}
		err = c.base(r, l)
		if err != nil {
			return fmt.Errorf("error filling company from %s: %w", f, err)
		}
		f, err = c.toJSON(dir)
		if err != nil {
			return err
		}
	}
	return nil
}

type baseTask struct {
	outDir  string
	lookups *lookups
	queues  []chan []string
	success chan struct{}
	errors  chan error
	bar     *progressbar.ProgressBar
}

func (t *baseTask) consumeShard(n int) {
	for r := range t.queues[n] {
		if err := addBase(t.lookups, t.outDir, r); err != nil {
			t.errors <- fmt.Errorf("error processing base cnpj %v: %w", r, err)
			continue
		}
		t.success <- struct{}{}
	}
}

func addBases(srcDir, outDir string, l *lookups) error {
	s, err := newSource(base, srcDir)
	if err != nil {
		return fmt.Errorf("error creating source for base cnpj: %w", err)
	}
	defer s.close()
	if s.totalLines == 0 {
		return nil
	}

	t := baseTask{
		outDir:  outDir,
		lookups: l,
		success: make(chan struct{}),
		errors:  make(chan error),
		bar:     progressbar.Default(s.totalLines),
	}
	t.bar.Describe("Adding base CNPJ")
	for i := 0; i < numOfShards; i++ {
		t.queues = append(t.queues, make(chan []string))
	}
	for i := 0; i < numOfShards; i++ {
		go t.consumeShard(i)
	}
	for _, r := range s.readers {
		go func(a *archivedCSV, q []chan []string, e chan error) {
			defer a.close()
			for {
				r, err := a.read()
				if err == io.EOF {
					break
				}
				if err != nil {
					e <- fmt.Errorf("error reading line %v: %w", r, err)
					break // do not proceed in case of errors.
				}
				s, err := shard(r[0])
				if err != nil {
					e <- fmt.Errorf("error getting shard for %s: %w", r[0], err)
					break // do not proceed in case of errors.
				}
				t.queues[s] <- r
			}
		}(r, t.queues, t.errors)
	}

	defer func() {
		for _, q := range t.queues {
			close(q)
		}
		close(t.success)
		close(t.errors)
	}()

	for {
		select {
		case err := <-t.errors:
			return err
		case <-t.success:
			t.bar.Add(1)
			if t.bar.IsFinished() {
				return nil
			}
		}
	}
}
