package tile

import (
	"fmt"
	"os"

	sm "github.com/flopp/go-staticmaps"
)

// getTileProviders gets the list of available OpenStreetMap tile providers
func getTileProviders() []string {
	tps := sm.GetTileProviders()
	ret := []string{}
	for _, t := range tps {
		ret = append(ret, t.Name)
	}
	return ret
}

// ListTileProvider lists all tile providers
func ListTileProvider() {
	tps := getTileProviders()
	fmt.Printf("Available tile providers\n")
	fmt.Printf("------------------------\n")
	for _, n := range tps {
		fmt.Println(n)
	}
	os.Exit(0)
}

// ValidateTileProvider is for checking cli args
func ValidateTileProvider(s string) bool {
	for _, tp := range getTileProviders() {
		if tp == s {
			return true
		}
	}
	return false
}

func ProviderByName(name string) *sm.TileProvider {
	for _, tp := range sm.GetTileProviders() {
		if tp.Name == name {
			return tp
		}
	}
	return &sm.TileProvider{}
}
