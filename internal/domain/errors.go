package domain

import "errors"

var (
	ErrStokHabis = errors.New("stok buku tidak tersedia")
	ErrStokPenuh = errors.New("stok sudah penuh")
	ErrBukuTidakDitemukan = errors.New("buku tidak ditemukan")
)