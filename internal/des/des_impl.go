package des

//go:generate go run ../../cmd/generate_data

import "encoding/binary"

type Des struct {
	subkeys [16]uint64
}

func New(key []uint8, mode_encrypt bool) *Des {
	des := &Des{}
	des.setKey(key, mode_encrypt)
	return des
}

func update_param(param *uint32, shift_left uint8) {
	shift_right := 28 - shift_left
	*param = (*param << shift_left) | ((*param >> shift_right) & 0xFFFFFFF0)
}

func (d *Des) setKey(key_bytes []uint8, mode_encrypt bool) {
	key := binary.LittleEndian.Uint64(key_bytes)

	param := map_u64(key, key_permutation_table)
	param_c := u64_get_lo32(param)
	param_d := u64_get_hi32(param)

	for i, shift_left := range key_rnd_shifts {
		var subkey_idx int
		if mode_encrypt {
			subkey_idx = i
		} else {
			subkey_idx = 15 - i
		}

		update_param(&param_c, shift_left)
		update_param(&param_d, shift_left)

		d.subkeys[subkey_idx] = map_u64(make_u64(param_d, param_c), key_compression)
	}
}

func (d *Des) TransformBlock(data uint64) uint64 {
	state := des_ip(data)
	for _, key := range d.subkeys {
		state = des_crypt_proc(state, key)
	}
	state = swap_u64_side(state)
	return des_ip_inv(state)
}

func (d *Des) TransformBytes(data []uint8) bool {
	if len(data)%8 != 0 {
		return false
	}

	for i := 0; i < len(data); i += 8 {
		value := binary.LittleEndian.Uint64(data[i : i+8])
		value = d.TransformBlock(value)
		binary.LittleEndian.PutUint64(data[i:i+8], value)
	}

	return true
}

func des_ip(data uint64) uint64 {
	return map_u64(data, ip)
}

func des_ip_inv(data uint64) uint64 {
	return map_u64(data, ip_inv)
}

func des_crypt_proc(state uint64, key uint64) uint64 {
	state_hi32 := u64_get_hi32(state)
	state_lo32 := u64_get_lo32(state)

	state = map_u64(make_u64(state_hi32, state_hi32), key_expansion)
	state ^= key

	next_lo32 := sbox_transform(state)
	next_lo32 = map_u32_bits(next_lo32, p_box)
	next_lo32 ^= state_lo32

	return make_u64(next_lo32, state_hi32)
}

func sbox_transform(state uint64) uint32 {
	result := uint32(0)
	for i, large_state_shift_v := range large_state_shifts {
		sbox_idx := (state >> uint64(large_state_shift_v)) & 0b111111
		result = (result << 4) | uint32(sboxes[i][sbox_idx])
	}
	return result
}
