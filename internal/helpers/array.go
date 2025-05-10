package helpers

import "fmt"

func CompareArrays(arr1, arr2 []string) bool {
	arr2Map := make(map[string]bool)

	// Mengisi peta dengan elemen-elemen dari arr2
	for _, item := range arr2 {
		arr2Map[item] = true
	}

	// Memeriksa apakah semua elemen dari arr1 terdapat dalam arr2
	for _, item := range arr1 {
		if !arr2Map[item] {
			return false // Ada elemen di arr1 yang tidak ada di arr2, mengembalikan false
		}
	}

	return true // Semua elemen dari arr1 terdapat dalam arr2, mengembalikan true
}

func IsStringInArray(str string, arr []string) bool {

	for _, s := range arr {
		fmt.Println("ini arr nya =", s)
		fmt.Println("ini str nya =", str)
		if s == str {
			return true
		}
	}
	return false
}
