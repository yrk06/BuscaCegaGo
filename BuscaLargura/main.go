package main

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type Tabuleiro [3][3]int

type Node struct {
	tab    Tabuleiro
	parent *Node
}

func gerarTabuleiro(v int) Tabuleiro {
	// 0 é espaço vazio
	number := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	rand.Seed(time.Now().UnixNano())
	//rand.Shuffle(len(number), func(i, j int) { number[i], number[j] = number[j], number[i] })
	var table Tabuleiro
	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			table[y][x] = number[3*y+x]
		}
	}
	for b := 0; b < v; b++ {
		child := gerarEstadosFilhos(table)
		table = child[rand.Intn(len(child))]
	}
	return table
}

func completo(tabuleiro Tabuleiro) bool {
	return tabuleiro == Tabuleiro{{0, 1, 2}, {3, 4, 5}, {6, 7, 8}}
}

func hashEstado(tabuleiro Tabuleiro) string {
	str := ""
	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			str += string(tabuleiro[y][x])
		}
	}

	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

func gerarEstadosFilhos(tabuleiro Tabuleiro) []Tabuleiro {
	X, Y := 0, 0
	for idy, y := range tabuleiro {
		for idx, x := range y {
			if x == 0 {
				X, Y = idx, idy
			}
		}
	}
	var estados_novos []Tabuleiro

	for x := -1; x <= 1; x += 2 {
		idx := X + x
		if idx < 0 || idx > 2 {
			continue
		}
		estado := tabuleiro
		estado[Y][X] = estado[Y][idx]
		estado[Y][idx] = 0
		estados_novos = append(estados_novos, estado)
	}

	for y := -1; y <= 1; y += 2 {
		idx := Y + y
		if idx < 0 || idx > 2 {
			continue
		}
		estado := tabuleiro
		estado[Y][X] = estado[idx][X]
		estado[idx][X] = 0
		estados_novos = append(estados_novos, estado)
	}

	return estados_novos
}

func buscaEmLargura(tabuleiro Tabuleiro, b int) []Tabuleiro {

	var visitados []string
	fila := []Node{Node{tabuleiro, nil}}
	for len(fila) > 0 {
		current := fila[0]
		hash := hashEstado(current.tab)
		visitados = append(visitados, hash)

		if completo(current.tab) {
			var passos []Tabuleiro
			passos = append(passos, current.tab)
			next := current.parent
			for next != nil {
				passos = append(passos, (*next).tab)
				next = (*next).parent
			}
			for i, j := 0, len(passos)-1; i < j; i, j = i+1, j-1 {
				passos[i], passos[j] = passos[j], passos[i]
			}
			return passos

		}

		for _, estado := range gerarEstadosFilhos(current.tab) {
			hash := hashEstado(estado)
			visitado := false
			for _, hsh := range visitados {
				if hsh == hash {
					visitado = true
					break
				}
			}
			if !visitado {
				fila = append(fila, Node{estado, &current})
			}

		}

		fila = fila[1:]
	}

	return nil
}

func buscaEmProfundidadeIterativa(tabuleiro Tabuleiro, maxb int) []Tabuleiro {
	for b := 1; b < maxb; b++ {
		resultado := buscaEmLargura(tabuleiro, b)
		if len(resultado) > 0 {
			return resultado
		}
	}
	return nil
}

func main() {

	var tcompleto Tabuleiro
	var base Tabuleiro
	var depth int

	for idx, arg := range os.Args {
		if arg == "-i" {
			for b := 0; b < 9; b++ {
				val, err := strconv.Atoi(os.Args[idx+b+1])
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				tcompleto[b/3][b%3] = val
			}
			fmt.Printf("Using user providaded state")
		}
		if arg == "-d" {
			d, err := strconv.Atoi(os.Args[idx+1])
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			depth = d
			fmt.Printf("Using max depth: %d\n", depth)
		}
	}

	if tcompleto == base {
		if depth != 0 {
			tcompleto = gerarTabuleiro(depth)
		} else {
			tcompleto = gerarTabuleiro(150)
		}

	}
	fmt.Println("Starting state:")
	for _, row := range tcompleto {
		fmt.Println(row)
	}
	fmt.Println("-solve-")
	for _, state := range buscaEmLargura(tcompleto, -1) {
		for _, row := range state {
			fmt.Println(row)
		}
		fmt.Println("-------")
	}

}
