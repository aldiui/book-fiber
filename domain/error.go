package domain

import "errors"

var BookNotFound = errors.New("data buku tidak ditemukan")
var JournalNotFound = errors.New("data jurnal tidak ditemukan")
var CustomerNotFound = errors.New("data customer tidak ditemukan")
var UserNotFound = errors.New("data user tidak ditemukan")
var BookStockNotAvailable = errors.New("stok buku tidak tersedia")
