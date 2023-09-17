package main

import (
	"fmt"
	"os"
)

func print_array_u8(f *os.File, name string, value []uint8) {
	fmt.Fprintf(f, "var %s = []uint8 {", name)
	for i := 0; i < len(value); i++ {
		if i%16 == 0 {
			fmt.Fprintf(f, "\n\t")
		}
		fmt.Fprintf(f, "0x%02x, ", value[i])
	}
	fmt.Fprintf(f, "\n}\n\n")
}

func print_array_u64(f *os.File, name string, value []uint64) {
	fmt.Fprintf(f, "var %s = []uint64 {", name)
	for i := 0; i < len(value); i++ {
		if i%4 == 0 {
			fmt.Fprintf(f, "\n\t")
		}
		fmt.Fprintf(f, "0x%016x, ", value[i])
	}
	fmt.Fprintf(f, "\n}\n\n")
}

func print_sbox(f *os.File, name string, value [8][64]uint8) {
	fmt.Fprintf(f, "var %s = [8][64]uint8 {\n", name)
	for i := 0; i < 8; i++ {
		fmt.Fprintf(f, "\t{ ")
		for j := 0; j < 64; j++ {
			if j != 0 {
				fmt.Fprintf(f, ", ")
			}
			fmt.Fprintf(f, "%d", value[i][j])
		}
		fmt.Fprintf(f, " },\n")
	}
	fmt.Fprintf(f, "}\n\n")
}

func main() {
	f, err := os.Create("./data.go")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fmt.Fprintf(f, "// Code generated by cmd/generate_data.go; DO NOT EDIT.\n\n")
	fmt.Fprintf(f, "package des\n\n")
	fmt.Fprintf(f, "// des_shift_table_cache is used for fast 32-bit system lookup.\n")
	print_array_u8(f, "key_rnd_shifts", []uint8{1, 1, 2, 2, 2, 2, 2, 2, 1, 2, 2, 2, 2, 2, 2, 1})
	print_array_u8(f, "large_state_shifts", []uint8{26, 20, 14, 8, 58, 52, 46, 40})
	print_sbox(f, "sboxes", make_des_sbox())
	print_array_u8(f, "p_box", data_sub_1([]uint8{
		16, 7, 20, 21, 29, 12, 28, 17, 1, 15, 23, 26, 5, 18, 31, 10, 2, 8, 24, 14, 32, 27, 3, 9, 19,
		13, 30, 6, 22, 11, 4, 25,
	}))
	print_array_u8(f, "ip", data_sub_1([]uint8{
		58, 50, 42, 34, 26, 18, 10, 2, 60, 52, 44, 36, 28, 20, 12, 4, 62, 54, 46, 38, 30, 22, 14, 6,
		64, 56, 48, 40, 32, 24, 16, 8, 57, 49, 41, 33, 25, 17, 9, 1, 59, 51, 43, 35, 27, 19, 11, 3, 61,
		53, 45, 37, 29, 21, 13, 5, 63, 55, 47, 39, 31, 23, 15, 7,
	}))
	print_array_u8(f, "ip_inv", data_sub_1([]uint8{
		40, 8, 48, 16, 56, 24, 64, 32, 39, 7, 47, 15, 55, 23, 63, 31, 38, 6, 46, 14, 54, 22, 62, 30,
		37, 5, 45, 13, 53, 21, 61, 29, 36, 4, 44, 12, 52, 20, 60, 28, 35, 3, 43, 11, 51, 19, 59, 27,
		34, 2, 42, 10, 50, 18, 58, 26, 33, 1, 41, 9, 49, 17, 57, 25,
	}))
	print_array_u8(f, "key_permutation_table", data_sub_1([]uint8{
		//key_param_c
		57, 49, 41, 33, 25, 17, 9, 1, 58, 50, 42, 34, 26, 18, 10, 2, 59, 51, 43, 35, 27, 19, 11, 3, 60,
		52, 44, 36,
		//key_param_d
		63, 55, 47, 39, 31, 23, 15, 7, 62, 54, 46, 38, 30, 22, 14, 6, 61, 53, 45, 37, 29, 21, 13, 5, 28,
		20, 12, 4,
	}))
	print_array_u8(f, "key_compression", make_key_compression())
	print_array_u8(f, "key_expansion", data_sub_1([]uint8{
		32, 1, 2, 3, 4, 5, 4, 5, 6, 7, 8, 9, 8, 9, 10, 11, 12, 13, 12, 13, 14, 15, 16, 17,
		16, 17, 18, 19, 20, 21, 20, 21, 22, 23, 24, 25, 24, 25, 26, 27, 28, 29, 28, 29, 30, 31, 32, 1,
	}))

	// For some reason, golang does not allow optimisation of negative shlifts.
	print_array_u64(f, "des_shift_table_cache", make_des_shift_table_cache())
}

func data_sub_1(data []uint8) []uint8 {
	return data_add(data, 0xff)
}

func data_add(data []uint8, delta uint8) []uint8 {
	result := make([]uint8, len(data))
	for i := 0; i < len(data); i++ {
		result[i] = data[i] + delta
	}
	return result
}

func make_des_shift_table_cache() []uint64 {
	des_shift_table_cache := [64]uint64{}
	for i := 0; i < 32; i++ {
		des_shift_table_cache[i] = 1 << (31 - i)
		des_shift_table_cache[i+32] = 1 << (31 - i + 32)
	}
	return des_shift_table_cache[:]
}

func make_key_compression() []uint8 {
	result := []uint8{
		// part 1
		14, 17, 11, 24, 1, 5, 3, 28, 15, 6, 21, 10, 23, 19, 12, 4, 26, 8, 16, 7, 27, 20, 13,
		2, //
		// part 2
		41, 52, 31, 37, 47, 55, 30, 40, 51, 45, 33, 48, 44, 49, 39, 56, 34, 53, 46, 42, 50, 36, 29,
		32,
	}

	for i := 0; i < 24; i++ {
		result[i] -= 1
	}
	for i := 24; i < 48; i++ {
		result[i] = result[i] - 28 + 32
	}

	return result
}

func make_des_sbox() [8][64]uint8 {
	return [8][64]uint8{
		{
			14, 0, 4, 15, 13, 7, 1, 4, 2, 14, 15, 2, 11, 13, 8, 1, 3, 10, 10, 6, 6, 12, 12, 11, 5, 9,
			9, 5, 0, 3, 7, 8, 4, 15, 1, 12, 14, 8, 8, 2, 13, 4, 6, 9, 2, 1, 11, 7, 15, 5, 12, 11, 9, 3,
			7, 14, 3, 10, 10, 0, 5, 6, 0, 13,
		},
		{
			15, 3, 1, 13, 8, 4, 14, 7, 6, 15, 11, 2, 3, 8, 4, 15, 9, 12, 7, 0, 2, 1, 13, 10, 12, 6, 0,
			9, 5, 11, 10, 5, 0, 13, 14, 8, 7, 10, 11, 1, 10, 3, 4, 15, 13, 4, 1, 2, 5, 11, 8, 6, 12, 7,
			6, 12, 9, 0, 3, 5, 2, 14, 15, 9,
		},
		{
			10, 13, 0, 7, 9, 0, 14, 9, 6, 3, 3, 4, 15, 6, 5, 10, 1, 2, 13, 8, 12, 5, 7, 14, 11, 12, 4,
			11, 2, 15, 8, 1, 13, 1, 6, 10, 4, 13, 9, 0, 8, 6, 15, 9, 3, 8, 0, 7, 11, 4, 1, 15, 2, 14,
			12, 3, 5, 11, 10, 5, 14, 2, 7, 12,
		},
		{
			7, 13, 13, 8, 14, 11, 3, 5, 0, 6, 6, 15, 9, 0, 10, 3, 1, 4, 2, 7, 8, 2, 5, 12, 11, 1, 12,
			10, 4, 14, 15, 9, 10, 3, 6, 15, 9, 0, 0, 6, 12, 10, 11, 10, 7, 13, 13, 8, 15, 9, 1, 4, 3,
			5, 14, 11, 5, 12, 2, 7, 8, 2, 4, 14,
		},
		{
			2, 14, 12, 11, 4, 2, 1, 12, 7, 4, 10, 7, 11, 13, 6, 1, 8, 5, 5, 0, 3, 15, 15, 10, 13, 3, 0,
			9, 14, 8, 9, 6, 4, 11, 2, 8, 1, 12, 11, 7, 10, 1, 13, 14, 7, 2, 8, 13, 15, 6, 9, 15, 12, 0,
			5, 9, 6, 10, 3, 4, 0, 5, 14, 3,
		},
		{
			12, 10, 1, 15, 10, 4, 15, 2, 9, 7, 2, 12, 6, 9, 8, 5, 0, 6, 13, 1, 3, 13, 4, 14, 14, 0, 7,
			11, 5, 3, 11, 8, 9, 4, 14, 3, 15, 2, 5, 12, 2, 9, 8, 5, 12, 15, 3, 10, 7, 11, 0, 14, 4, 1,
			10, 7, 1, 6, 13, 0, 11, 8, 6, 13,
		},
		{
			4, 13, 11, 0, 2, 11, 14, 7, 15, 4, 0, 9, 8, 1, 13, 10, 3, 14, 12, 3, 9, 5, 7, 12, 5, 2, 10,
			15, 6, 8, 1, 6, 1, 6, 4, 11, 11, 13, 13, 8, 12, 1, 3, 4, 7, 10, 14, 7, 10, 9, 15, 5, 6, 0,
			8, 15, 0, 14, 5, 2, 9, 3, 2, 12,
		},
		{
			13, 1, 2, 15, 8, 13, 4, 8, 6, 10, 15, 3, 11, 7, 1, 4, 10, 12, 9, 5, 3, 6, 14, 11, 5, 0, 0,
			14, 12, 9, 7, 2, 7, 2, 11, 1, 4, 14, 1, 7, 9, 4, 12, 10, 14, 8, 2, 13, 0, 15, 6, 12, 10, 9,
			13, 0, 15, 3, 3, 5, 5, 6, 8, 11,
		},
	}
}
