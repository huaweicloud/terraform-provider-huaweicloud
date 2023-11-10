// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Based on CRYPTOGAMS code with the following comment:
// # ====================================================================
// # Written by Andy Polyakov <appro@openssl.org> for the OpenSSL
// # project. The module is, however, dual licensed under OpenSSL and
// # CRYPTOGAMS licenses depending on where you obtain it. For further
// # details see http://www.openssl.org/~appro/cryptogams/.
// # ====================================================================

// Code for the perl script that generates the ppc64 assembler
// can be found in the cryptogams repository at the link below. It is based on
// the original from openssl.

// https://github.com/dot-asm/cryptogams/commit/a60f5b50ed908e91

// The differences in this and the original implementation are
// due to the calling conventions and initialization of constants.

//go:build gc && !purego
// +build gc,!purego

#include "textflag.h"

#define OUT  R3
#define INP  R4
#define LEN  R5
#define KEY  R6
#define CNT  R7
#define TMP  R15

#define CONSTBASE  R16
#define BLOCKS R17

DATA consts<>+0x00(SB)/8, $0x3320646e61707865
DATA consts<>+0x08(SB)/8, $0x6b20657479622d32
DATA consts<>+0x10(SB)/8, $0x0000000000000001
DATA consts<>+0x18(SB)/8, $0x0000000000000000
DATA consts<>+0x20(SB)/8, $0x0000000000000004
DATA consts<>+0x28(SB)/8, $0x0000000000000000
DATA consts<>+0x30(SB)/8, $0x0a0b08090e0f0c0d
DATA consts<>+0x38(SB)/8, $0x0203000106070405
DATA consts<>+0x40(SB)/8, $0x090a0b080d0e0f0c
DATA consts<>+0x48(SB)/8, $0x0102030005060704
DATA consts<>+0x50(SB)/8, $0x6170786561707865
DATA consts<>+0x58(SB)/8, $0x6170786561707865
DATA consts<>+0x60(SB)/8, $0x3320646e3320646e
DATA consts<>+0x68(SB)/8, $0x3320646e3320646e
DATA consts<>+0x70(SB)/8, $0x79622d3279622d32
DATA consts<>+0x78(SB)/8, $0x79622d3279622d32
DATA consts<>+0x80(SB)/8, $0x6b2065746b206574
DATA consts<>+0x88(SB)/8, $0x6b2065746b206574
DATA consts<>+0x90(SB)/8, $0x0000000100000000
DATA consts<>+0x98(SB)/8, $0x0000000300000002
GLOBL consts<>(SB), RODATA, $0xa0

//func chaCha20_ctr32_vsx(out, inp *byte, len int, key *[8]uint32, counter *uint32)
TEXT ·chaCha20_ctr32_vsx(SB),NOSPLIT,$64-40
	MOVD out+0(FP), OUT
	MOVD inp+8(FP), INP
	MOVD len+16(FP), LEN
	MOVD key+24(FP), KEY
	MOVD counter+32(FP), CNT

	// Addressing for constants
	MOVD $consts<>+0x00(SB), CONSTBASE
	MOVD $16, R8
	MOVD $32, R9
	MOVD $48, R10
	MOVD $64, R11
	SRD $6, LEN, BLOCKS
	// V16
	LXVW4X (CONSTBASE)(R0), VS48
	ADD $80,CONSTBASE

	// Load key into V17,V18
	LXVW4X (KEY)(R0), VS49
	LXVW4X (KEY)(R8), VS50

	// Load CNT, NONCE into V19
	LXVW4X (CNT)(R0), VS51

	// Clear V27
	VXOR V27, V27, V27

	// V28
	LXVW4X (CONSTBASE)(R11), VS60

	// splat slot from V19 -> V26
	VSPLTW $0, V19, V26

	VSLDOI $4, V19, V27, V19
	VSLDOI $12, V27, V19, V19

	VADDUWM V26, V28, V26

	MOVD $10, R14
	MOVD R14, CTR

loop_outer_vsx:
	// V0, V1, V2, V3
	LXVW4X (R0)(CONSTBASE), VS32
	LXVW4X (R8)(CONSTBASE), VS33
	LXVW4X (R9)(CONSTBASE), VS34
	LXVW4X (R10)(CONSTBASE), VS35

	// splat values from V17, V18 into V4-V11
	VSPLTW $0, V17, V4
	VSPLTW $1, V17, V5
	VSPLTW $2, V17, V6
	VSPLTW $3, V17, V7
	VSPLTW $0, V18, V8
	VSPLTW $1, V18, V9
	VSPLTW $2, V18, V10
	VSPLTW $3, V18, V11

	// VOR
	VOR V26, V26, V12

	// splat values from V19 -> V13, V14, V15
	VSPLTW $1, V19, V13
	VSPLTW $2, V19, V14
	VSPLTW $3, V19, V15

	// splat   const values
	VSPLTISW $-16, V27
	VSPLTISW $12, V28
	VSPLTISW $8, V29
	VSPLTISW $7, V30

loop_vsx:
	VADDUWM V0, V4, V0
	VADDUWM V1, V5, V1
	VADDUWM V2, V6, V2
	VADDUWM V3, V7, V3

	VXOR V12, V0, V12
	VXOR V13, V1, V13
	VXOR V14, V2, V14
	VXOR V15, V3, V15

	VRLW V12, V27, V12
	VRLW V13, V27, V13
	VRLW V14, V27, V14
	VRLW V15, V27, V15

	VADDUWM V8, V12, V8
	VADDUWM V9, V13, V9
	VADDUWM V10, V14, V10
	VADDUWM V11, V15, V11

	VXOR V4, V8, V4
	VXOR V5, V9, V5
	VXOR V6, V10, V6
	VXOR V7, V11, V7

	VRLW V4, V28, V4
	VRLW V5, V28, V5
	VRLW V6, V28, V6
	VRLW V7, V28, V7

	VADDUWM V0, V4, V0
	VADDUWM V1, V5, V1
	VADDUWM V2, V6, V2
	VADDUWM V3, V7, V3

	VXOR V12, V0, V12
	VXOR V13, V1, V13
	VXOR V14, V2, V14
	VXOR V15, V3, V15

	VRLW V12, V29, V12
	VRLW V13, V29, V13
	VRLW V14, V29, V14
	VRLW V15, V29, V15

	VADDUWM V8, V12, V8
	VADDUWM V9, V13, V9
	VADDUWM V10, V14, V10
	VADDUWM V11, V15, V11

	VXOR V4, V8, V4
	VXOR V5, V9, V5
	VXOR V6, V10, V6
	VXOR V7, V11, V7

	VRLW V4, V30, V4
	VRLW V5, V30, V5
	VRLW V6, V30, V6
	VRLW V7, V30, V7

	VADDUWM V0, V5, V0
	VADDUWM V1, V6, V1
	VADDUWM V2, V7, V2
	VADDUWM V3, V4, V3

	VXOR V15, V0, V15
	VXOR V12, V1, V12
	VXOR V13, V2, V13
	VXOR V14, V3, V14

	VRLW V15, V27, V15
	VRLW V12, V27, V12
	VRLW V13, V27, V13
	VRLW V14, V27, V14

	VADDUWM V10, V15, V10
	VADDUWM V11, V12, V11
	VADDUWM V8, V13, V8
	VADDUWM V9, V14, V9

	VXOR V5, V10, V5
	VXOR V6, V11, V6
	VXOR V7, V8, V7
	VXOR V4, V9, V4

	VRLW V5, V28, V5
	VRLW V6, V28, V6
	VRLW V7, V28, V7
	VRLW V4, V28, V4

	VADDUWM V0, V5, V0
	VADDUWM V1, V6, V1
	VADDUWM V2, V7, V2
	VADDUWM V3, V4, V3

	VXOR V15, V0, V15
	VXOR V12, V1, V12
	VXOR V13, V2, V13
	VXOR V14, V3, V14

	VRLW V15, V29, V15
	VRLW V12, V29, V12
	VRLW V13, V29, V13
	VRLW V14, V29, V14

	VADDUWM V10, V15, V10
	VADDUWM V11, V12, V11
	VADDUWM V8, V13, V8
	VADDUWM V9, V14, V9

	VXOR V5, V10, V5
	VXOR V6, V11, V6
	VXOR V7, V8, V7
	VXOR V4, V9, V4

	VRLW V5, V30, V5
	VRLW V6, V30, V6
	VRLW V7, V30, V7
	VRLW V4, V30, V4
	BC   16, LT, loop_vsx

	VADDUWM V12, V26, V12

	WORD $0x13600F8C		// VMRGEW V0, V1, V27
	WORD $0x13821F8C		// VMRGEW V2, V3, V28

	WORD $0x10000E8C		// VMRGOW V0, V1, V0
	WORD $0x10421E8C		// VMRGOW V2, V3, V2

	WORD $0x13A42F8C		// VMRGEW V4, V5, V29
	WORD $0x13C63F8C		// VMRGEW V6, V7, V30

	XXPERMDI VS32, VS34, $0, VS33
	XXPERMDI VS32, VS34, $3, VS35
	XXPERMDI VS59, VS60, $0, VS32
	XXPERMDI VS59, VS60, $3, VS34

	WORD $0x10842E8C		// VMRGOW V4, V5, V4
	WORD $0x10C63E8C		// VMRGOW V6, V7, V6

	WORD $0x13684F8C		// VMRGEW V8, V9, V27
	WORD $0x138A5F8C		// VMRGEW V10, V11, V28

	XXPERMDI VS36, VS38, $0, VS37
	XXPERMDI VS36, VS38, $3, VS39
	XXPERMDI VS61, VS62, $0, VS36
	XXPERMDI VS61, VS62, $3, VS38

	WORD $0x11084E8C		// VMRGOW V8, V9, V8
	WORD $0x114A5E8C		// VMRGOW V10, V11, V10

	WORD $0x13AC6F8C		// VMRGEW V12, V13, V29
	WORD $0x13CE7F8C		// VMRGEW V14, V15, V30

	XXPERMDI VS40, VS42, $0, VS41
	XXPERMDI VS40, VS42, $3, VS43
	XXPERMDI VS59, VS60, $0, VS40
	XXPERMDI VS59, VS60, $3, VS42

	WORD $0x118C6E8C		// VMRGOW V12, V13, V12
	WORD $0x11CE7E8C		// VMRGOW V14, V15, V14

	VSPLTISW $4, V27
	VADDUWM V26, V27, V26

	XXPERMDI VS44, VS46, $0, VS45
	XXPERMDI VS44, VS46, $3, VS47
	XXPERMDI VS61, VS62, $0, VS44
	XXPERMDI VS61, VS62, $3, VS46

	VADDUWM V0, V16, V0
	VADDUWM V4, V17, V4
	VADDUWM V8, V18, V8
	VADDUWM V12, V19, V12

	CMPU LEN, $64
	BLT tail_vsx

	// Bottom of loop
	LXVW4X (INP)(R0), VS59
	LXVW4X (INP)(R8), VS60
	LXVW4X (INP)(R9), VS61
	LXVW4X (INP)(R10), VS62

	VXOR V27, V0, V27
	VXOR V28, V4, V28
	VXOR V29, V8, V29
	VXOR V30, V12, V30

	STXVW4X VS59, (OUT)(R0)
	STXVW4X VS60, (OUT)(R8)
	ADD     $64, INP
	STXVW4X VS61, (OUT)(R9)
	ADD     $-64, LEN
	STXVW4X VS62, (OUT)(R10)
	ADD     $64, OUT
	BEQ     done_vsx

	VADDUWM V1, V16, V0
	VADDUWM V5, V17, V4
	VADDUWM V9, V18, V8
	VADDUWM V13, V19, V12

	CMPU  LEN, $64
	BLT   tail_vsx

	LXVW4X (INP)(R0), VS59
	LXVW4X (INP)(R8), VS60
	LXVW4X (INP)(R9), VS61
	LXVW4X (INP)(R10), VS62
	VXOR   V27, V0, V27

	VXOR V28, V4, V28
	VXOR V29, V8, V29
	VXOR V30, V12, V30

	STXVW4X VS59, (OUT)(R0)
	STXVW4X VS60, (OUT)(R8)
	ADD     $64, INP
	STXVW4X VS61, (OUT)(R9)
	ADD     $-64, LEN
	STXVW4X VS62, (OUT)(V10)
	ADD     $64, OUT
	BEQ     done_vsx

	VADDUWM V2, V16, V0
	VADDUWM V6, V17, V4
	VADDUWM V10, V18, V8
	VADDUWM V14, V19, V12

	CMPU LEN, $64
	BLT  tail_vsx

	LXVW4X (INP)(R0), VS59
	LXVW4X (INP)(R8), VS60
	LXVW4X (INP)(R9), VS61
	LXVW4X (INP)(R10), VS62

	VXOR V27, V0, V27
	VXOR V28, V4, V28
	VXOR V29, V8, V29
	VXOR V30, V12, V30

	STXVW4X VS59, (OUT)(R0)
	STXVW4X VS60, (OUT)(R8)
	ADD     $64, INP
	STXVW4X VS61, (OUT)(R9)
	ADD     $-64, LEN
	STXVW4X VS62, (OUT)(R10)
	ADD     $64, OUT
	BEQ     done_vsx

	VADDUWM V3, V16, V0
	VADDUWM V7, V17, V4
	VADDUWM V11, V18, V8
	VADDUWM V15, V19, V12

	CMPU  LEN, $64
	BLT   tail_vsx

	LXVW4X (INP)(R0), VS59
	LXVW4X (INP)(R8), VS60
	LXVW4X (INP)(R9), VS61
	LXVW4X (INP)(R10), VS62

	VXOR V27, V0, V27
	VXOR V28, V4, V28
	VXOR V29, V8, V29
	VXOR V30, V12, V30

	STXVW4X VS59, (OUT)(R0)
	STXVW4X VS60, (OUT)(R8)
	ADD     $64, INP
	STXVW4X VS61, (OUT)(R9)
	ADD     $-64, LEN
	STXVW4X VS62, (OUT)(R10)
	ADD     $64, OUT

	MOVD $10, R14
	MOVD R14, CTR
	BNE  loop_outer_vsx

done_vsx:
	// Increment counter by number of 64 byte blocks
	MOVD (CNT), R14
	ADD  BLOCKS, R14
	MOVD R14, (CNT)
	RET

tail_vsx:
	ADD  $32, R1, R11
	MOVD LEN, CTR

	// Save values on stack to copy from
	STXVW4X VS32, (R11)(R0)
	STXVW4X VS36, (R11)(R8)
	STXVW4X VS40, (R11)(R9)
	STXVW4X VS44, (R11)(R10)
	ADD $-1, R11, R12
	ADD $-1, INP
	ADD $-1, OUT

looptail_vsx:
	// Copying the result to OUT
	// in bytes.
	MOVBZU 1(R12), KEY
	MOVBZU 1(INP), TMP
	XOR    KEY, TMP, KEY
	MOVBU  KEY, 1(OUT)
	BC     16, LT, looptail_vsx

	// Clear the stack values
	STXVW4X VS48, (R11)(R0)
	STXVW4X VS48, (R11)(R8)
	STXVW4X VS48, (R11)(R9)
	STXVW4X VS48, (R11)(R10)
	BR      done_vsx
