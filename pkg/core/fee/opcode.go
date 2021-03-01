package fee

import (
	"github.com/nspcc-dev/neo-go/pkg/vm/opcode"
)

// Opcode returns the deployment coefficients of specified opcodes.
func Opcode(base int64, opcodes ...opcode.Opcode) int64 {
	var result int64
	for _, op := range opcodes {
		result += coefficients[op]
	}
	return result * base
}

var coefficients = map[opcode.Opcode]int64{
	opcode.PUSHINT8:     1 << 0,
	opcode.PUSHINT16:    1 << 0,
	opcode.PUSHINT32:    1 << 0,
	opcode.PUSHINT64:    1 << 0,
	opcode.PUSHINT128:   1 << 2,
	opcode.PUSHINT256:   1 << 2,
	opcode.PUSHA:        1 << 2,
	opcode.PUSHNULL:     1 << 0,
	opcode.PUSHDATA1:    1 << 3,
	opcode.PUSHDATA2:    1 << 9,
	opcode.PUSHDATA4:    1 << 12,
	opcode.PUSHM1:       1 << 0,
	opcode.PUSH0:        1 << 0,
	opcode.PUSH1:        1 << 0,
	opcode.PUSH2:        1 << 0,
	opcode.PUSH3:        1 << 0,
	opcode.PUSH4:        1 << 0,
	opcode.PUSH5:        1 << 0,
	opcode.PUSH6:        1 << 0,
	opcode.PUSH7:        1 << 0,
	opcode.PUSH8:        1 << 0,
	opcode.PUSH9:        1 << 0,
	opcode.PUSH10:       1 << 0,
	opcode.PUSH11:       1 << 0,
	opcode.PUSH12:       1 << 0,
	opcode.PUSH13:       1 << 0,
	opcode.PUSH14:       1 << 0,
	opcode.PUSH15:       1 << 0,
	opcode.PUSH16:       1 << 0,
	opcode.NOP:          1 << 0,
	opcode.JMP:          1 << 1,
	opcode.JMPL:         1 << 1,
	opcode.JMPIF:        1 << 1,
	opcode.JMPIFL:       1 << 1,
	opcode.JMPIFNOT:     1 << 1,
	opcode.JMPIFNOTL:    1 << 1,
	opcode.JMPEQ:        1 << 1,
	opcode.JMPEQL:       1 << 1,
	opcode.JMPNE:        1 << 1,
	opcode.JMPNEL:       1 << 1,
	opcode.JMPGT:        1 << 1,
	opcode.JMPGTL:       1 << 1,
	opcode.JMPGE:        1 << 1,
	opcode.JMPGEL:       1 << 1,
	opcode.JMPLT:        1 << 1,
	opcode.JMPLTL:       1 << 1,
	opcode.JMPLE:        1 << 1,
	opcode.JMPLEL:       1 << 1,
	opcode.CALL:         1 << 9,
	opcode.CALLL:        1 << 9,
	opcode.CALLA:        1 << 9,
	opcode.ABORT:        0,
	opcode.ASSERT:       1 << 0,
	opcode.THROW:        1 << 9,
	opcode.TRY:          1 << 2,
	opcode.TRYL:         1 << 2,
	opcode.ENDTRY:       1 << 2,
	opcode.ENDTRYL:      1 << 2,
	opcode.ENDFINALLY:   1 << 2,
	opcode.RET:          0,
	opcode.SYSCALL:      0,
	opcode.DEPTH:        1 << 1,
	opcode.DROP:         1 << 1,
	opcode.NIP:          1 << 1,
	opcode.XDROP:        1 << 4,
	opcode.CLEAR:        1 << 4,
	opcode.DUP:          1 << 1,
	opcode.OVER:         1 << 1,
	opcode.PICK:         1 << 1,
	opcode.TUCK:         1 << 1,
	opcode.SWAP:         1 << 1,
	opcode.ROT:          1 << 1,
	opcode.ROLL:         1 << 4,
	opcode.REVERSE3:     1 << 1,
	opcode.REVERSE4:     1 << 1,
	opcode.REVERSEN:     1 << 4,
	opcode.INITSSLOT:    1 << 4,
	opcode.INITSLOT:     1 << 6,
	opcode.LDSFLD0:      1 << 1,
	opcode.LDSFLD1:      1 << 1,
	opcode.LDSFLD2:      1 << 1,
	opcode.LDSFLD3:      1 << 1,
	opcode.LDSFLD4:      1 << 1,
	opcode.LDSFLD5:      1 << 1,
	opcode.LDSFLD6:      1 << 1,
	opcode.LDSFLD:       1 << 1,
	opcode.STSFLD0:      1 << 1,
	opcode.STSFLD1:      1 << 1,
	opcode.STSFLD2:      1 << 1,
	opcode.STSFLD3:      1 << 1,
	opcode.STSFLD4:      1 << 1,
	opcode.STSFLD5:      1 << 1,
	opcode.STSFLD6:      1 << 1,
	opcode.STSFLD:       1 << 1,
	opcode.LDLOC0:       1 << 1,
	opcode.LDLOC1:       1 << 1,
	opcode.LDLOC2:       1 << 1,
	opcode.LDLOC3:       1 << 1,
	opcode.LDLOC4:       1 << 1,
	opcode.LDLOC5:       1 << 1,
	opcode.LDLOC6:       1 << 1,
	opcode.LDLOC:        1 << 1,
	opcode.STLOC0:       1 << 1,
	opcode.STLOC1:       1 << 1,
	opcode.STLOC2:       1 << 1,
	opcode.STLOC3:       1 << 1,
	opcode.STLOC4:       1 << 1,
	opcode.STLOC5:       1 << 1,
	opcode.STLOC6:       1 << 1,
	opcode.STLOC:        1 << 1,
	opcode.LDARG0:       1 << 1,
	opcode.LDARG1:       1 << 1,
	opcode.LDARG2:       1 << 1,
	opcode.LDARG3:       1 << 1,
	opcode.LDARG4:       1 << 1,
	opcode.LDARG5:       1 << 1,
	opcode.LDARG6:       1 << 1,
	opcode.LDARG:        1 << 1,
	opcode.STARG0:       1 << 1,
	opcode.STARG1:       1 << 1,
	opcode.STARG2:       1 << 1,
	opcode.STARG3:       1 << 1,
	opcode.STARG4:       1 << 1,
	opcode.STARG5:       1 << 1,
	opcode.STARG6:       1 << 1,
	opcode.STARG:        1 << 1,
	opcode.NEWBUFFER:    1 << 8,
	opcode.MEMCPY:       1 << 11,
	opcode.CAT:          1 << 11,
	opcode.SUBSTR:       1 << 11,
	opcode.LEFT:         1 << 11,
	opcode.RIGHT:        1 << 11,
	opcode.INVERT:       1 << 2,
	opcode.AND:          1 << 3,
	opcode.OR:           1 << 3,
	opcode.XOR:          1 << 3,
	opcode.EQUAL:        1 << 5,
	opcode.NOTEQUAL:     1 << 5,
	opcode.SIGN:         1 << 2,
	opcode.ABS:          1 << 2,
	opcode.NEGATE:       1 << 2,
	opcode.INC:          1 << 2,
	opcode.DEC:          1 << 2,
	opcode.ADD:          1 << 3,
	opcode.SUB:          1 << 3,
	opcode.MUL:          1 << 3,
	opcode.DIV:          1 << 3,
	opcode.MOD:          1 << 3,
	opcode.POW:          1 << 6,
	opcode.SQRT:         1 << 11,
	opcode.SHL:          1 << 3,
	opcode.SHR:          1 << 3,
	opcode.NOT:          1 << 2,
	opcode.BOOLAND:      1 << 3,
	opcode.BOOLOR:       1 << 3,
	opcode.NZ:           1 << 2,
	opcode.NUMEQUAL:     1 << 3,
	opcode.NUMNOTEQUAL:  1 << 3,
	opcode.LT:           1 << 3,
	opcode.LTE:          1 << 3,
	opcode.GT:           1 << 3,
	opcode.GTE:          1 << 3,
	opcode.MIN:          1 << 3,
	opcode.MAX:          1 << 3,
	opcode.WITHIN:       1 << 3,
	opcode.PACK:         1 << 9,
	opcode.UNPACK:       1 << 9,
	opcode.NEWARRAY0:    1 << 4,
	opcode.NEWARRAY:     1 << 9,
	opcode.NEWARRAYT:    1 << 9,
	opcode.NEWSTRUCT0:   1 << 4,
	opcode.NEWSTRUCT:    1 << 9,
	opcode.NEWMAP:       1 << 3,
	opcode.SIZE:         1 << 2,
	opcode.HASKEY:       1 << 6,
	opcode.KEYS:         1 << 4,
	opcode.VALUES:       1 << 13,
	opcode.PICKITEM:     1 << 6,
	opcode.APPEND:       1 << 13,
	opcode.SETITEM:      1 << 13,
	opcode.REVERSEITEMS: 1 << 13,
	opcode.REMOVE:       1 << 4,
	opcode.CLEARITEMS:   1 << 4,
	opcode.POPITEM:      1 << 4,
	opcode.ISNULL:       1 << 1,
	opcode.ISTYPE:       1 << 1,
	opcode.CONVERT:      1 << 11,
}
