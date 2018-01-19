// Copyright 2018 The Golem Language Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
)

const (
	LOAD_NULL byte = iota
	LOAD_TRUE
	LOAD_FALSE
	LOAD_ZERO
	LOAD_ONE
	LOAD_NEG_ONE

	IMPORT_MODULE
	LOAD_BUILTIN
	LOAD_CONST
	LOAD_LOCAL
	STORE_LOCAL
	LOAD_CAPTURE
	STORE_CAPTURE

	JUMP
	JUMP_TRUE
	JUMP_FALSE

	EQ
	NE
	GT
	GTE
	LT
	LTE
	CMP
	HAS

	PLUS
	SUB
	MUL
	DIV

	REM
	BIT_AND
	BIT_OR
	BIT_XOR
	LEFT_SHIFT
	RIGHT_SHIFT

	NEGATE
	NOT
	COMPLEMENT

	NEW_FUNC
	FUNC_CAPTURE
	FUNC_LOCAL

	INVOKE
	GO
	RETURN
	DONE
	THROW

	NEW_STRUCT
	NEW_DICT
	NEW_LIST
	NEW_SET
	NEW_TUPLE

	GET_FIELD
	INIT_FIELD
	SET_FIELD
	INC_FIELD

	GET_INDEX
	SET_INDEX
	INC_INDEX

	SLICE
	SLICE_FROM
	SLICE_TO

	ITER
	ITER_NEXT
	ITER_GET

	CHECK_CAST
	CHECK_TUPLE

	POP
	DUP

	// These are temporary values created during compilation.
	// The interpreter will panic if it encounters them.
	BREAK    = 0xFD
	CONTINUE = 0xFE
)

func OpCodeSize(opc byte) int {

	switch opc {

	case IMPORT_MODULE, LOAD_BUILTIN, LOAD_CONST,
		LOAD_LOCAL, LOAD_CAPTURE, STORE_LOCAL, STORE_CAPTURE,
		JUMP, JUMP_TRUE, JUMP_FALSE, BREAK, CONTINUE,
		NEW_FUNC, FUNC_CAPTURE, FUNC_LOCAL, INVOKE, GO,
		NEW_STRUCT, GET_FIELD, INIT_FIELD, SET_FIELD, INC_FIELD,
		NEW_DICT, NEW_LIST, NEW_SET, NEW_TUPLE, CHECK_CAST, CHECK_TUPLE:

		return 3

	default:
		return 1
	}
}

func FmtOpcode(opcodes []byte, i int) string {

	switch opcodes[i] {

	case LOAD_NULL:
		return fmt.Sprintf("%d: LOAD_NULL\n", i)
	case LOAD_TRUE:
		return fmt.Sprintf("%d: LOAD_TRUE\n", i)
	case LOAD_FALSE:
		return fmt.Sprintf("%d: LOAD_FALSE\n", i)
	case LOAD_ZERO:
		return fmt.Sprintf("%d: LOAD_ZERO\n", i)
	case LOAD_ONE:
		return fmt.Sprintf("%d: LOAD_ONE\n", i)
	case LOAD_NEG_ONE:
		return fmt.Sprintf("%d: LOAD_NEG_ONE\n", i)

	case IMPORT_MODULE:
		return fmtIndex(opcodes, i, "IMPORT_MODULE")
	case LOAD_BUILTIN:
		return fmtIndex(opcodes, i, "LOAD_BUILTIN")
	case LOAD_CONST:
		return fmtIndex(opcodes, i, "LOAD_CONST")
	case LOAD_LOCAL:
		return fmtIndex(opcodes, i, "LOAD_LOCAL")
	case STORE_LOCAL:
		return fmtIndex(opcodes, i, "STORE_LOCAL")
	case LOAD_CAPTURE:
		return fmtIndex(opcodes, i, "LOAD_CAPTURE")
	case STORE_CAPTURE:
		return fmtIndex(opcodes, i, "STORE_CAPTURE")

	case JUMP:
		return fmtIndex(opcodes, i, "JUMP")
	case JUMP_TRUE:
		return fmtIndex(opcodes, i, "JUMP_TRUE")
	case JUMP_FALSE:
		return fmtIndex(opcodes, i, "JUMP_FALSE")

	case EQ:
		return fmt.Sprintf("%d: EQ\n", i)
	case NE:
		return fmt.Sprintf("%d: NE\n", i)
	case GT:
		return fmt.Sprintf("%d: GT\n", i)
	case GTE:
		return fmt.Sprintf("%d: GTE\n", i)
	case LT:
		return fmt.Sprintf("%d: LT\n", i)
	case LTE:
		return fmt.Sprintf("%d: LTE\n", i)
	case CMP:
		return fmt.Sprintf("%d: CMP\n", i)
	case HAS:
		return fmt.Sprintf("%d: HAS\n", i)

	case PLUS:
		return fmt.Sprintf("%d: PLUS\n", i)
	case SUB:
		return fmt.Sprintf("%d: SUB\n", i)
	case MUL:
		return fmt.Sprintf("%d: MUL\n", i)
	case DIV:
		return fmt.Sprintf("%d: DIV\n", i)

	case REM:
		return fmt.Sprintf("%d: REM\n", i)
	case BIT_AND:
		return fmt.Sprintf("%d: BIT_AND\n", i)
	case BIT_OR:
		return fmt.Sprintf("%d: BIT_OR\n", i)
	case BIT_XOR:
		return fmt.Sprintf("%d: BIT_XOR\n", i)
	case LEFT_SHIFT:
		return fmt.Sprintf("%d: LEFT_SHIFT\n", i)
	case RIGHT_SHIFT:
		return fmt.Sprintf("%d: RIGHT_SHIFT\n", i)

	case NEGATE:
		return fmt.Sprintf("%d: NEGATE\n", i)
	case NOT:
		return fmt.Sprintf("%d: NOT\n", i)
	case COMPLEMENT:
		return fmt.Sprintf("%d: COMPLEMENT\n", i)

	case NEW_FUNC:
		return fmtIndex(opcodes, i, "NEW_FUNC")
	case FUNC_CAPTURE:
		return fmtIndex(opcodes, i, "FUNC_CAPTURE")
	case FUNC_LOCAL:
		return fmtIndex(opcodes, i, "FUNC_LOCAL")

	case INVOKE:
		return fmtIndex(opcodes, i, "INVOKE")
	case GO:
		return fmtIndex(opcodes, i, "GO")
	case RETURN:
		return fmt.Sprintf("%d: RETURN\n", i)
	case DONE:
		return fmt.Sprintf("%d: DONE\n", i)
	case THROW:
		return fmt.Sprintf("%d: THROW\n", i)

	case NEW_STRUCT:
		return fmtIndex(opcodes, i, "NEW_STRUCT")
	case GET_FIELD:
		return fmtIndex(opcodes, i, "GET_FIELD")
	case INIT_FIELD:
		return fmtIndex(opcodes, i, "INIT_FIELD")
	case SET_FIELD:
		return fmtIndex(opcodes, i, "SET_FIELD")
	case INC_FIELD:
		return fmtIndex(opcodes, i, "INC_FIELD")
	case NEW_DICT:
		return fmtIndex(opcodes, i, "NEW_DICT")
	case NEW_LIST:
		return fmtIndex(opcodes, i, "NEW_LIST")
	case NEW_SET:
		return fmtIndex(opcodes, i, "NEW_SET")
	case NEW_TUPLE:
		return fmtIndex(opcodes, i, "NEW_TUPLE")

	case GET_INDEX:
		return fmt.Sprintf("%d: GET_INDEX\n", i)
	case SET_INDEX:
		return fmt.Sprintf("%d: SET_INDEX\n", i)
	case INC_INDEX:
		return fmt.Sprintf("%d: INC_INDEX\n", i)

	case SLICE:
		return fmt.Sprintf("%d: SLICE\n", i)
	case SLICE_FROM:
		return fmt.Sprintf("%d: SLICE_FROM\n", i)
	case SLICE_TO:
		return fmt.Sprintf("%d: SLICE_TO\n", i)

	case ITER:
		return fmt.Sprintf("%d: ITER\n", i)
	case ITER_NEXT:
		return fmt.Sprintf("%d: ITER_NEXT\n", i)
	case ITER_GET:
		return fmt.Sprintf("%d: ITER_GET\n", i)

	case CHECK_CAST:
		return fmtIndex(opcodes, i, "CHECK_CAST")
	case CHECK_TUPLE:
		return fmtIndex(opcodes, i, "CHECK_TUPLE")

	case POP:
		return fmt.Sprintf("%d: POP\n", i)
	case DUP:
		return fmt.Sprintf("%d: DUP\n", i)

	case BREAK:
		return fmtIndex(opcodes, i, "BREAK")
	case CONTINUE:
		return fmtIndex(opcodes, i, "CONTINUE")

	default:
		panic(fmt.Sprintf("unreachable %d", opcodes[i]))
	}
}

func fmtIndex(opcodes []byte, i int, tag string) string {
	high := opcodes[i+1]
	low := opcodes[i+2]
	index := int(high)<<8 + int(low)
	return fmt.Sprintf("%d: %s %d %d (%d)\n", i, tag, high, low, index)
}
