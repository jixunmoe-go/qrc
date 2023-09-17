package des

func make_u64(hi32, lo32 uint32) uint64 {
	return uint64(hi32)<<32 | uint64(lo32)
}

func swap_u64_side(value uint64) uint64 {
	return (value >> 32) | (value << 32)
}

func u64_get_lo32(value uint64) uint32 {
	return uint32(value)
}

func u64_get_hi32(value uint64) uint32 {
	return uint32(value >> 32)
}

func get_u64_by_shift_idx(value uint8) uint64 {
	return des_shift_table_cache[value&0x3f]
}

func map_bit(result *uint64, src uint64, check, set uint8) {
	if (get_u64_by_shift_idx(check) & src) != 0 {
		*result |= get_u64_by_shift_idx(set)
	}
}

func map_u32_bits(src_value uint32, table []uint8) uint32 {
	result := uint64(0)
	for i, v := range table {
		map_bit(&result, uint64(src_value), v, uint8(i))
	}
	return uint32(result)
}

func map_u64(src_value uint64, table []uint8) uint64 {
	mid_idx := len(table) / 2

	table_lo32 := table[:mid_idx]
	table_hi32 := table[mid_idx:]

	lo32 := uint64(0)
	hi32 := uint64(0)

	for i, v := range table_lo32 {
		map_bit(&lo32, src_value, v, uint8(i))
	}
	for i, v := range table_hi32 {
		map_bit(&hi32, src_value, v, uint8(i))
	}

	return make_u64(uint32(hi32), uint32(lo32))
}
