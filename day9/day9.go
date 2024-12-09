package day9

import (
	"bufio"
	"context"
	"embed"
	"strconv"

	. "github.com/gverger/aoc2024/utils"
)

//go:embed input.txt
//go:embed sample.txt
var f embed.FS

type DiskMap []int

func (d DiskMap) FileID(idx int) int {
	Assert(idx%2 == 0, "FileID is odd")

	return idx / 2
}

type Input struct {
	DiskMap DiskMap
}

func ReadInput(filename string) Input {
	file := Must(f.Open(filename))
	defer file.Close()

	scanner := bufio.NewScanner(file)

	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	AssertNoErr(scanner.Err(), "reading input file")
	Assert(len(lines) == 1, "Should be one line, but found %d", len(lines))
	line := lines[0]

	diskmap := make(DiskMap, len(line))
	for i, c := range line {
		diskmap[i] = Must(strconv.Atoi(string(c)))
	}

	return Input{
		DiskMap: diskmap,
	}
}

type InputLoaded struct {
	Input Input
}

type SolutionFound struct {
	Part     int
	Solution int64
}

func compactIndividualChunks(_ context.Context, diskmap DiskMap) int {
	dmIdx := 0
	sum := 0

	blockID := 0
	for dmIdx < len(diskmap) && diskmap[dmIdx] > 0 {
		fileID := diskmap.FileID(dmIdx)

		for i := 0; i < diskmap[dmIdx]; i++ {
			sum += fileID * blockID
			// log.Info().Int("sum", sum).Interface("diskmap", diskmap).Msg("read")
			blockID++
		}

		dmIdx++
		if dmIdx >= len(diskmap) {
			break
		}

		for i := 0; i < diskmap[dmIdx]; i++ {
			lastIdx := len(diskmap) - 1
			if diskmap[lastIdx] == 0 {
				diskmap = diskmap[:len(diskmap)-2]
				lastIdx = len(diskmap) - 1
			}
			fileID := diskmap.FileID(lastIdx)

			diskmap[lastIdx]--
			sum += fileID * blockID
			blockID++
			// log.Info().Int("sum", sum).Interface("diskmap", diskmap).Msg("compact")
		}

		dmIdx++

	}

	return sum
}

func firstAvailableSpot(diskmap DiskMap, nbBlocks int) int {
	for i := 1; i < len(diskmap); i += 2 {
		if diskmap[i] >= nbBlocks {
			return i
		}
	}
	return -1
}

func compactWholeFiles(_ context.Context, diskmap DiskMap) int64 {
	var sum int64

	startingBlock := make([]int, len(diskmap))
	currentBlock := 0
	for i, v := range diskmap {
		startingBlock[i] = currentBlock
		currentBlock += v
	}

	for lastIdx := len(diskmap) - 1; lastIdx >= 0; lastIdx -= 2 {
		fileID := diskmap.FileID(lastIdx)
		nbBlocks := diskmap[lastIdx]

		blockID := startingBlock[lastIdx]
		slot := firstAvailableSpot(diskmap[:lastIdx], nbBlocks)
		if slot != -1 {
			diskmap[slot] -= nbBlocks
			blockID = startingBlock[slot]
			startingBlock[slot] += nbBlocks
		}

		for i := 0; i < nbBlocks; i++ {
			sum += int64(blockID * fileID)
			blockID++
		}
	}

	return sum
}

func Run(ctx context.Context, callback func(ctx context.Context, obj any)) {

	input := ReadInput("input.txt")
	callback(ctx, InputLoaded{Input: input})

	diskmap := make(DiskMap, len(input.DiskMap))
	copy(diskmap, input.DiskMap)
	callback(ctx, SolutionFound{Part: 1, Solution: int64(compactIndividualChunks(ctx, diskmap))})

	diskmap = make(DiskMap, len(input.DiskMap))
	copy(diskmap, input.DiskMap)
	callback(ctx, SolutionFound{Part: 2, Solution: compactWholeFiles(ctx, diskmap)})
}
