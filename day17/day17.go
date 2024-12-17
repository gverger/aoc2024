package day17

import (
	"bufio"
	"context"
	"embed"
	"strconv"
	"strings"

	utils "github.com/gverger/aoc2024/utils"
	"github.com/phuslu/log"
)

//go:embed *.txt
var f embed.FS

type Computer struct {
	A int
	B int
	C int

	Instructions []int
	Pointer      int

	Out string
}

func (c Computer) LiteralOperand() int {
	return c.Instructions[c.Pointer+1]
}

func (c Computer) ComboOperand() int {
	return c.OperandValue(c.LiteralOperand())
}

func (c *Computer) Run(program ...int) {
	c.Instructions = program
	c.Pointer = 0
	c.Out = ""

	for c.Pointer < len(program) {
		c.run()
	}
}

func (c *Computer) run() {
	c.OpCode()(c)
	c.Pointer += 2
}

func (c Computer) OperandValue(operand int) int {
	utils.Assert(operand >= 0 && operand <= 6, "operand invalid")
	if operand <= 3 {
		return operand
	}
	if operand == 4 {
		return c.A
	}
	if operand == 5 {
		return c.B
	}
	if operand == 6 {
		return c.C
	}

	log.Fatal().Msg("cannot reach here")
	return -1
}

type instruction func(c *Computer)

func adv(c *Computer) {
	c.A = c.A >> c.ComboOperand()
}

func bxl(c *Computer) {
	c.B = c.B ^ c.LiteralOperand()
}

func bst(c *Computer) {
	c.B = c.ComboOperand() & 7
}

func jnz(c *Computer) {
	if c.A == 0 {
		return
	}
	c.Pointer = c.LiteralOperand() - 2
}

func bxc(c *Computer) {
	c.B = c.B ^ c.C
}

func out(c *Computer) {
	sep := ""
	if len(c.Out) > 0 {
		sep = ","
	}
	c.Out += sep + strconv.Itoa(c.ComboOperand()&7)
}

func bdv(c *Computer) {
	c.B = c.A >> c.ComboOperand()
}

func cdv(c *Computer) {
	c.C = c.A >> c.ComboOperand()
}

func (c Computer) OpCode() instruction {
	switch c.Instructions[c.Pointer] {
	case 0:
		return adv
	case 1:
		return bxl
	case 2:
		return bst
	case 3:
		return jnz
	case 4:
		return bxc
	case 5:
		return out
	case 6:
		return bdv
	case 7:
		return cdv
	}

	return nil
}

type Input struct {
	A       int
	B       int
	C       int
	Program []int
}

func lastField(line string) string {
	fields := strings.Fields(line)
	return fields[len(fields)-1]
}

func ReadInput(filename string) Input {
	file := utils.Must(f.Open(filename))
	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Scan()
	a := utils.Must(strconv.Atoi(lastField(scanner.Text())))
	scanner.Scan()
	b := utils.Must(strconv.Atoi(lastField(scanner.Text())))
	scanner.Scan()
	c := utils.Must(strconv.Atoi(lastField(scanner.Text())))
	scanner.Scan()
	scanner.Scan()
	program := utils.MapTo(
		strings.FieldsFunc(lastField(scanner.Text()), func(r rune) bool { return r == ',' }),
		func(s string) int { return utils.Must(strconv.Atoi(s)) },
	)

	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	utils.AssertNoErr(scanner.Err(), "reading input file")

	return Input{
		A:       a,
		B:       b,
		C:       c,
		Program: program,
	}
}

type InputLoaded struct {
	Input Input
}

type SolutionFound struct {
	Part     int
	Solution string
}

func nextStep(c Computer, digit int) (int, bool) {
	initA := c.A
	wanted := strings.Join(
		utils.MapTo(c.Instructions[digit:],
			func(i int) string { return strconv.Itoa(i) }), ",")

	for v := 0; v < 8; v++ {
		c.A = initA*8 + v
		c.Run(c.Instructions...)

		if c.Out == wanted {
			a := initA*8 + v
			if digit == 0 {
				return a, true
			}
			c.A = a
			if newA, ok := nextStep(c, digit-1); ok {
				return newA, true
			}
		}
	}

	return 0, false
}

func unroll(c Computer) (int, bool) {
	c.A = 0
	return nextStep(c, len(c.Instructions)-1)
}

func Run(ctx context.Context, callback func(ctx context.Context, obj any)) {
	filename := "input.txt"
	input := ReadInput(filename)

	callback(ctx, InputLoaded{Input: input})

	c := Computer{
		A: input.A,
		B: input.B,
		C: input.C,
	}

	c.Run(input.Program...)

	callback(ctx, SolutionFound{Part: 1, Solution: c.Out})

	c = Computer{
		A:            input.A,
		B:            input.B,
		C:            input.C,
		Instructions: input.Program,
	}
	n, ok := unroll(c)
	utils.Assert(ok, "Not found")

	callback(ctx, SolutionFound{Part: 2, Solution: strconv.Itoa(n)})
}
