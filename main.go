package main

import (
	"html/template"
	"math/rand"
	"net"
	"net/http"
)

const size = 100
const probability = 0.4

type point [2]int

var nbd = []point{
	{-1,-1},{0,-1},{1,-1},
	{-1,0},{1,0},
	{-1,1},{0,1},{1,1},
}

func (p point) add(q point) point {
	return point{p[0]+q[0], p[1]+q[1]}
}

var lsn net.Listener
var t *template.Template
var board [][]bool
var alive map[point]bool
var changed map[point]bool

const tbl = `<table>{{range .}}<tr>{{range .}}<td{{if .}} class="on"{{end}}/>{{end}}</tr>{{end}}</table>`

var handlers = map[string]func(http.ResponseWriter, *http.Request){
	"/":       root,
	"/quit":   quit,
	"/random": random,
	"/step":   step,
}

func main() {
	t = template.New("table")
	t.Parse(tbl)
	for k, v := range handlers {
		http.HandleFunc(k, v)
	}
	lsn, _ = net.Listen("tcp", ":8080")
	http.Serve(lsn, nil)
}

func quit(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Shutting down..."))
	lsn.Close()
}

func random(w http.ResponseWriter, r *http.Request) {
	board = make([][]bool, size)
	alive = make(map[point]bool)
	changed = make(map[point]bool)
	for i := range board {
		r := make([]bool, size)
		board[i] = r
		for j := range r {
			p := [2]int{i,j}
			changed[p] = true
			if rand.Float32() < probability {
				r[j] = true
				alive[p] = true
			}
		}
	}
	t.Execute(w, board)
}

func live(me bool, count int) bool {
	if count == 3 || count == 4 && me {
		return true
	}
	return false
	//Any live cell with fewer than two live neighbours dies, as if caused by under-population.
	//Any live cell with two or three live neighbours lives on to the next generation.
	//Any live cell with more than three live neighbours dies, as if by over-population.
	//Any dead cell with exactly three live neighbours becomes a live cell, as if by reproduction.
}

func step(w http.ResponseWriter, r *http.Request) {
	create := make(map[point]bool)
	kill := make(map[point]bool)
	chg := make(map[point]bool)
	
	for p := range changed {
		count := 0
		for _, n := range nbd {
			if alive[p.add(n)] {
				count++
			}
		}
		me := alive[p]
		if me {
			count++
		}
		l := live(me, count)
		if me != l {
			chg[p] = true
			for _, n := range nbd {
				chg[p.add(n)] = true
			}
			x,y := p[0], p[1]
			if x < 0 || y < 0 || x >= len(board) || y >= len(board[x]) {
				continue
			}
			if l {
				create[p] = true
				board[x][y] = true
			} else {
				kill[p] = true
				board[x][y] = false
			}
		}
	}
	
	for p := range kill {
		delete(alive, p)
	}
	for p := range create {
		alive[p] = true
	}
	
	changed = chg
//	count := make([][]int, size)
//	for i, row := range board {
//		count[i] = make([]int, size)
//		for j := range row {
//			c := 0
//			if j > 0 && board[i][j-1] {
//				c++
//			}
//			if board[i][j] {
//				c++
//			}
//			if j < size-1 && board[i][j+1] {
//				c++
//			}
//			count[i][j] = c
//		}
//	}

//	for i, row := range count {
//		for j := range row {
//			c := 0
//			if i > 0 {
//				c += count[i-1][j]
//			}
//			c += count[i][j]
//			if i < size-1 {
//				c += count[i+1][j]
//			}
//			board[i][j] = live(board[i][j], c)
//		}
//	}

	t.Execute(w, board)
}

func root(w http.ResponseWriter, r *http.Request) {
	if len(r.RequestURI) == 1 {
		r.RequestURI = "/life.html"
	}
	http.ServeFile(w, r, r.RequestURI[1:])
}
