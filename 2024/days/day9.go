package days

import (
	"fmt"
	"strconv"
	"time"

	"github.com/blister/adventofcode/2024/color"
)

type BlockStore struct {
	store []*Block
}

type Block struct {
	file_id    int
	buffer     bool
	buffer_len int
	block_len  int
}

func ReadDisk(data string) *BlockStore {

	store := &BlockStore{store: make([]*Block, 0)}

	//fmt.Println("INPUT", data)

	var file_id int
	for i := 0; i < len(data); i++ {
		file_blocks, err := strconv.Atoi(string(data[i]))
		if err != nil {
			panic(err)
		}

		var bufblocks int
		if i < len(data)-1 {
			buf, err := strconv.Atoi(string(data[i+1]))
			if err != nil {
				panic(err)
			}
			bufblocks = buf
			i++
		} else {
			bufblocks = 0
		}

		for k := 0; k < file_blocks; k++ {
			store.store = append(store.store, &Block{
				file_id:   file_id,
				buffer:    false,
				block_len: file_blocks,
			})

			//fmt.Print(strconv.Itoa(file_id))
		}
		file_id++
		for k := 0; k < bufblocks; k++ {
			buffer := &Block{
				file_id:    0,
				buffer:     true,
				buffer_len: bufblocks,
			}
			store.store = append(store.store, buffer)
			//fmt.Print(".")
		}
	}

	return store
}

func (b *BlockStore) Compress2(verbose bool) {
	for end_idx := len(b.store) - 1; end_idx > 0; {
		if b.store[end_idx].buffer == true {
			end_idx--
			continue
		} else {
			//fmt.Println("Trying", b.store[end_idx].file_id, b.store[end_idx].block_len)
			start_idx := 0
			for {
				// fmt.Println("\tSEARCH", b.store[end_idx].block_len, b.store[end_idx].file_id, start_idx, "buf_len", b.store[start_idx].buffer_len)

				if start_idx > end_idx {
					end_idx = end_idx - b.store[end_idx].block_len
					break
				}
				if b.store[start_idx].buffer == false {
					start_idx += b.store[start_idx].block_len
					continue
				} else if b.store[start_idx].buffer_len >= b.store[end_idx].block_len {
					// fmt.Println("else", b.store[start_idx].buffer_len, b.store[end_idx].block_len)
					var moved int = 0
					buffer_size := b.store[start_idx].buffer_len
					move_size := b.store[end_idx].block_len
					for {

						// fmt.Println("moved", moved, "buffer", buffer_size, "move", move_size)

						if moved < move_size {
							b.store[start_idx] = b.store[end_idx]
							b.store[end_idx] = &Block{
								file_id:    0,
								buffer:     true,
								buffer_len: move_size,
							}
							if verbose {
								fmt.Println("COMPRESS", b.Render(start_idx, end_idx))
							}
							moved++
							start_idx++
							end_idx--
						} else if moved < buffer_size {
							// fmt.Println("changing buffer size", buffer_size, buffer_size-move_size)
							b.store[start_idx].buffer_len = buffer_size - move_size
							moved++
							start_idx++
							continue
						} else {
							// fmt.Println("finished", moved, buffer_size, move_size)
							break
						}
					}
					break
				} else { // buffer < block_len
					start_idx++ //= b.store[start_idx].block_len
					continue
				}
			}
			if verbose {
				fmt.Println("compress", b.Render(start_idx, end_idx))
			}
		}
	}
}

func (b *BlockStore) Compress(verbose bool) {
	var end_idx int = len(b.store) - 1
	for start_idx := 0; start_idx < end_idx; {
		if b.store[start_idx].buffer == false {
			start_idx++
			continue
		} else {
			for {
				if b.store[end_idx].buffer == false {
					b.store[start_idx] = b.store[end_idx]
					b.store[end_idx] = &Block{
						file_id: 0,
						buffer:  true,
					}
					end_idx--
					break
				} else {
					end_idx--
					break
				}
			}
		}
		if verbose {
			fmt.Println("COMPRESS", b.Render(start_idx, end_idx))
		}
		// fmt.Println("COMPRESS", b.Render(start_idx, end_idx),
		// 	start_idx, "=", b.store[start_idx].buffer,
		// 	end_idx, "=", b.store[end_idx].buffer,
		// )
		//fmt.Println("COMPRESS", b.Render(start_idx, end_idx))
	}
}

func (b *BlockStore) CheckSum(report *Report) int {
	var checkSum int = 0
	report.debug = append(report.debug, "CHECKSUMS")
	for i, v := range b.store {
		if v.buffer {
			continue
		}
		checkSum += i * v.file_id
		report.debug = append(
			report.debug,
			fmt.Sprintf("%d * %d = %d", i, v.file_id, i*v.file_id),
		)
	}

	return checkSum
}

func (b *BlockStore) Render(start_idx int, end_idx int) string {
	var out string
	var use_color bool = false
	if end_idx > 0 {
		use_color = true
	}
	for i, v := range b.store {
		if use_color && i == start_idx {
			out += color.Cyan
		} else if use_color && i == end_idx {
			out += color.Green
		}
		if v.buffer {
			out += "."
		} else {
			out += strconv.Itoa(v.file_id)
		}

		if use_color {
			out += color.Reset
		}
	}
	return out
}

func Day9a(verbose bool, test bool, input string) Report {
	var report = Report{
		day:      "9a",
		solution: 0,
		start:    time.Now(),
	}
	report.correct = false
	report.stop = time.Now()

	var path string = "days/inputs/day9.txt"
	if test {
		path = "days/inputs/day9_test.txt"
	}
	if len(input) > 0 {
		path = "days/inputs/" + input
	}

	data, err := ReadFile(path)
	if err != nil {
		PrintError(err)
		report.solution = 0
		report.correct = false
		report.stop = time.Now()
		report.debug = append(report.debug, err.Error())
		return report
	}

	//fmt.Println(data)
	bs := ReadDisk(data)
	report.debug = append(report.debug, "INPUT:")
	report.debug = append(report.debug, bs.Render(0, 0))
	if verbose {
		fmt.Println("bs", len(bs.store))
	}
	bs.Compress(verbose)
	cs := bs.CheckSum(&report)

	report.solution = cs

	//bs.Compress()
	//checkSum := bs.CheckSum()

	// fmt.Println("CHECKSUM = ", checkSum)
	// report.solution = checkSum

	report.debug = append(report.debug, "COMPRESSED:")
	report.debug = append(report.debug, bs.Render(0, 0))

	report.correct = false
	report.stop = time.Now()

	return report
}

func Day9b(verbose bool, test bool, input string) Report {
	var report = Report{
		day:      "9b",
		solution: 0,
		start:    time.Now(),
	}
	report.correct = false
	report.stop = time.Now()

	var path string = "days/inputs/day9.txt"
	if test {
		path = "days/inputs/day9_test.txt"
	}
	if len(input) > 0 {
		path = "days/inputs/" + input
	}

	data, err := ReadFile(path)
	if err != nil {
		PrintError(err)
		report.solution = 0
		report.correct = false
		report.stop = time.Now()
		report.debug = append(report.debug, err.Error())
		return report
	}

	//fmt.Println(data)
	bs := ReadDisk(data)
	report.debug = append(report.debug, "INPUT:")
	report.debug = append(report.debug, bs.Render(0, 0))
	if verbose {
		fmt.Println("bs", len(bs.store))
	}
	bs.Compress2(verbose)
	cs := bs.CheckSum(&report)

	report.solution = cs

	//bs.Compress()
	//checkSum := bs.CheckSum()

	// fmt.Println("CHECKSUM = ", checkSum)
	// report.solution = checkSum

	report.debug = append(report.debug, "COMPRESSED:")
	report.debug = append(report.debug, bs.Render(0, 0))

	report.correct = false
	report.stop = time.Now()

	return report
}
