package download

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestFederalRevenueGetURLs(t *testing.T) {
	tmp := t.TempDir()
	ts := httpTestServer(t, "dados-publicos-cnpj.html")
	defer ts.Close()

	t.Run("returns download urls", func(t *testing.T) {
		got, err := federalRevenueGetURLs(ts.Client(), ts.URL, tmp)
		if err != nil {
			t.Errorf("expected to run withour errors, got: %v:", err)
		}
		expected := []string{
			"http://200.152.38.155/CNPJ/F.K03200$W.SIMPLES.CSV.D10710.zip",
			"http://200.152.38.155/CNPJ/F.K03200$Z.D10710.CNAECSV.zip",
			"http://200.152.38.155/CNPJ/F.K03200$Z.D10710.MOTICSV.zip",
			"http://200.152.38.155/CNPJ/F.K03200$Z.D10710.MUNICCSV.zip",
			"http://200.152.38.155/CNPJ/F.K03200$Z.D10710.NATJUCSV.zip",
			"http://200.152.38.155/CNPJ/F.K03200$Z.D10710.PAISCSV.zip",
			"http://200.152.38.155/CNPJ/F.K03200$Z.D10710.QUALSCSV.zip",
			"http://200.152.38.155/CNPJ/K3241.K03200Y0.D10710.EMPRECSV.zip",
			"http://200.152.38.155/CNPJ/K3241.K03200Y0.D10710.ESTABELE.zip",
			"http://200.152.38.155/CNPJ/K3241.K03200Y0.D10710.SOCIOCSV.zip",
			"http://200.152.38.155/CNPJ/K3241.K03200Y1.D10710.EMPRECSV.zip",
			"http://200.152.38.155/CNPJ/K3241.K03200Y1.D10710.ESTABELE.zip",
			"http://200.152.38.155/CNPJ/K3241.K03200Y1.D10710.SOCIOCSV.zip",
			"http://200.152.38.155/CNPJ/K3241.K03200Y2.D10710.EMPRECSV.zip",
			"http://200.152.38.155/CNPJ/K3241.K03200Y2.D10710.ESTABELE.zip",
			"http://200.152.38.155/CNPJ/K3241.K03200Y2.D10710.SOCIOCSV.zip",
			"http://200.152.38.155/CNPJ/K3241.K03200Y3.D10710.EMPRECSV.zip",
			"http://200.152.38.155/CNPJ/K3241.K03200Y3.D10710.ESTABELE.zip",
			"http://200.152.38.155/CNPJ/K3241.K03200Y3.D10710.SOCIOCSV.zip",
			"http://200.152.38.155/CNPJ/K3241.K03200Y4.D10710.EMPRECSV.zip",
			"http://200.152.38.155/CNPJ/K3241.K03200Y4.D10710.ESTABELE.zip",
			"http://200.152.38.155/CNPJ/K3241.K03200Y4.D10710.SOCIOCSV.zip",
			"http://200.152.38.155/CNPJ/K3241.K03200Y5.D10710.EMPRECSV.zip",
			"http://200.152.38.155/CNPJ/K3241.K03200Y5.D10710.ESTABELE.zip",
			"http://200.152.38.155/CNPJ/K3241.K03200Y5.D10710.SOCIOCSV.zip",
			"http://200.152.38.155/CNPJ/K3241.K03200Y6.D10710.EMPRECSV.zip",
			"http://200.152.38.155/CNPJ/K3241.K03200Y6.D10710.ESTABELE.zip",
			"http://200.152.38.155/CNPJ/K3241.K03200Y6.D10710.SOCIOCSV.zip",
			"http://200.152.38.155/CNPJ/K3241.K03200Y7.D10710.EMPRECSV.zip",
			"http://200.152.38.155/CNPJ/K3241.K03200Y7.D10710.ESTABELE.zip",
			"http://200.152.38.155/CNPJ/K3241.K03200Y7.D10710.SOCIOCSV.zip",
			"http://200.152.38.155/CNPJ/K3241.K03200Y8.D10710.EMPRECSV.zip",
			"http://200.152.38.155/CNPJ/K3241.K03200Y8.D10710.ESTABELE.zip",
			"http://200.152.38.155/CNPJ/K3241.K03200Y8.D10710.SOCIOCSV.zip",
			"http://200.152.38.155/CNPJ/K3241.K03200Y9.D10710.EMPRECSV.zip",
			"http://200.152.38.155/CNPJ/K3241.K03200Y9.D10710.ESTABELE.zip",
			"http://200.152.38.155/CNPJ/K3241.K03200Y9.D10710.SOCIOCSV.zip",
		}
		assertArraysHaveSameItems(t, got, expected)
	})

	t.Run("saves updated at date", func(t *testing.T) {
		_, err := federalRevenueGetURLs(ts.Client(), ts.URL, tmp)
		if err != nil {
			t.Errorf("expected to run withour errors, got: %v:", err)
		}
		pth := filepath.Join(tmp, federalRevenueUpdatedAt)
		got, err := ioutil.ReadFile(pth)
		if err != nil {
			t.Errorf("expected no error reading %s, updatedAt %s", pth, err)
		}
		expected := "2021-07-16"
		if string(got) != expected {
			t.Errorf("expected updated at to be %s, got %s", expected, string(got))
		}
	})
}
