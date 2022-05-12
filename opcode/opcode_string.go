// Code generated by "stringer -type=Opcode -output=opcode_string.go"; DO NOT EDIT.

package opcode

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Invalid-0]
	_ = x[Add-1]
	_ = x[Adds-2]
	_ = x[Addx-3]
	_ = x[And-4]
	_ = x[Andc-5]
	_ = x[Band-6]
	_ = x[Bcc-7]
	_ = x[Bclr-8]
	_ = x[Bcs-9]
	_ = x[Beq-10]
	_ = x[Bge-11]
	_ = x[Bgt-12]
	_ = x[Bhi-13]
	_ = x[Biand-14]
	_ = x[Bild-15]
	_ = x[Bior-16]
	_ = x[Bist-17]
	_ = x[Bixor-18]
	_ = x[Bld-19]
	_ = x[Ble-20]
	_ = x[Bls-21]
	_ = x[Blt-22]
	_ = x[Bmi-23]
	_ = x[Bne-24]
	_ = x[Bnot-25]
	_ = x[Bor-26]
	_ = x[Bpl-27]
	_ = x[Bra-28]
	_ = x[Brn-29]
	_ = x[Bset-30]
	_ = x[Bsr-31]
	_ = x[Bst-32]
	_ = x[Btst-33]
	_ = x[Bvc-34]
	_ = x[Bvs-35]
	_ = x[Bxor-36]
	_ = x[Clrmac-37]
	_ = x[Cmp-38]
	_ = x[Daa-39]
	_ = x[Das-40]
	_ = x[Dec-41]
	_ = x[Divxs-42]
	_ = x[Divxu-43]
	_ = x[Eepmov-44]
	_ = x[Exts-45]
	_ = x[Extu-46]
	_ = x[Inc-47]
	_ = x[Jmp-48]
	_ = x[Jsr-49]
	_ = x[Ldc-50]
	_ = x[Ldm-51]
	_ = x[Ldmac-52]
	_ = x[Mac-53]
	_ = x[Mov-54]
	_ = x[Movfpe-55]
	_ = x[Movtpe-56]
	_ = x[Mulxs-57]
	_ = x[Mulxu-58]
	_ = x[Neg-59]
	_ = x[Nop-60]
	_ = x[Not-61]
	_ = x[Or-62]
	_ = x[Orc-63]
	_ = x[Pop-64]
	_ = x[Push-65]
	_ = x[Rotl-66]
	_ = x[Rotr-67]
	_ = x[Rotxl-68]
	_ = x[Rotxr-69]
	_ = x[Rte-70]
	_ = x[Rts-71]
	_ = x[Shal-72]
	_ = x[Shar-73]
	_ = x[Shll-74]
	_ = x[Shlr-75]
	_ = x[Sleep-76]
	_ = x[Stc-77]
	_ = x[Stm-78]
	_ = x[Stmac-79]
	_ = x[Sub-80]
	_ = x[Subs-81]
	_ = x[Subx-82]
	_ = x[Tas-83]
	_ = x[Trapa-84]
	_ = x[Xor-85]
	_ = x[Xorc-86]
}

const _Opcode_name = "InvalidAddAddsAddxAndAndcBandBccBclrBcsBeqBgeBgtBhiBiandBildBiorBistBixorBldBleBlsBltBmiBneBnotBorBplBraBrnBsetBsrBstBtstBvcBvsBxorClrmacCmpDaaDasDecDivxsDivxuEepmovExtsExtuIncJmpJsrLdcLdmLdmacMacMovMovfpeMovtpeMulxsMulxuNegNopNotOrOrcPopPushRotlRotrRotxlRotxrRteRtsShalSharShllShlrSleepStcStmStmacSubSubsSubxTasTrapaXorXorc"

var _Opcode_index = [...]uint16{0, 7, 10, 14, 18, 21, 25, 29, 32, 36, 39, 42, 45, 48, 51, 56, 60, 64, 68, 73, 76, 79, 82, 85, 88, 91, 95, 98, 101, 104, 107, 111, 114, 117, 121, 124, 127, 131, 137, 140, 143, 146, 149, 154, 159, 165, 169, 173, 176, 179, 182, 185, 188, 193, 196, 199, 205, 211, 216, 221, 224, 227, 230, 232, 235, 238, 242, 246, 250, 255, 260, 263, 266, 270, 274, 278, 282, 287, 290, 293, 298, 301, 305, 309, 312, 317, 320, 324}

func (i Opcode) String() string {
	if i < 0 || i >= Opcode(len(_Opcode_index)-1) {
		return "Opcode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Opcode_name[_Opcode_index[i]:_Opcode_index[i+1]]
}
