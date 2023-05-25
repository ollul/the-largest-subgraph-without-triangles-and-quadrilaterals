package main

import (
        "container/list"
        "fmt"
        "io"
        "math/rand"
        "os"
        "sort"
        "strconv"
        "strings"
        "time"
)

var w1, w2 int32

var distance [][]uint16
var incidence [][]uint16
var edges []edge

type point struct {
        x, y uint16
}

type edge struct {
        v1, v2 int
}

func taxi(v1, v2 point) (res uint16) {
        if v1.x > v2.x {
                res += v1.x - v2.x
        } else {
                res += v2.x - v1.x
        }
        if v1.y > v2.y {
                res += v1.y - v2.y
        } else {
                res += v2.y - v1.y
        }
        return
}

type result struct {
        incidence [][]bool
        cost      int
}

var resChan chan *result

func processGraph(num int, thr int) {
        edgesc := make([]edge, len(edges))
        copy(edgesc, edges)
        lists := make([]map[int]bool, num)
        for {
                r := rand.New(rand.NewSource(time.Now().UnixNano()))
                incidencec := make([][]bool, num)
                for i := 0; i < num; i++ {
                        incidencec[i] = make([]bool, num)
                        lists[i] = make(map[int]bool)
                }
                k := len(edgesc)
                s := r.Int31n(w1) + 1
                for i := 0; i < int(s); i++ {
                        sort.Slice(edgesc[i:k/int(s)*(i+1)], func(i, j int) bool {
                                return distance[edgesc[i].v1][edgesc[i].v2] > distance[edgesc[j].v1][edgesc[j].v2]
                        })
                }
                for i := 0; i < int(r.Int31n(w2)); i++ {
                        w1 := r.Int31n(int32(k))
                        w2 := r.Int31n(int32(k))
                        edgesc[w1], edgesc[w2] = edgesc[w2], edgesc[w1]
                }
                edgesl := list.New()
                for _, k := range edgesc {
                        edgesl.PushBack(k)
                }
                cost := 0
                for {
                        var e edge
                        if edgesl.Front() == nil {
                                break
                        }
                        e = edgesl.Front().Value.(edge)
                        edgesl.Remove(edgesl.Front())
                        incidencec[e.v1][e.v2] = true
                        incidencec[e.v2][e.v1] = true
                        lists[e.v1][e.v2] = true
                        lists[e.v2][e.v1] = true
                        bad := false
                        for k := range lists[e.v1] {
                                if k != e.v2 && lists[k][e.v2] {
                                        bad = true
                                        break
                                }
                        }

                out:
                        for v1 := range lists[e.v1] {
                                for v2 := range lists[v1] {
                                        if e.v1 != v1 && e.v1 != v2 && e.v2 != v1 && e.v2 != v2 && lists[v2][e.v2] {
                                                bad = true
                                                break out
                                        }
                                }
                        }

                        if bad {
                                incidencec[e.v1][e.v2] = false
                                incidencec[e.v2][e.v1] = false
                                delete(lists[e.v1], e.v2)
                                delete(lists[e.v2], e.v1)
                        } else {
                                cost += int(distance[e.v1][e.v2])
                        }
                }
                resChan <- &result{
                        incidence: incidencec,
                        cost:      cost,
                }
        }
}

func main() {
        if len(os.Args) < 5 {
                fmt.Printf("Usage: %s vert num_threads w1 w2\n", os.Args[0])
                return
        }
        num, err := strconv.Atoi(os.Args[1])
        if err != nil {
                return
        }
        nt, err := strconv.Atoi(os.Args[2])
        if err != nil {
                return
        }
        w1i, err := strconv.Atoi(os.Args[3])
        if err != nil {
                return
        }
        w2i, err := strconv.Atoi(os.Args[3])
        if err != nil {
                return
        }
        w1 = int32(w1i)
        w2 = int32(w2i)
        vertices_file, _ := os.Open(fmt.Sprintf("../Taxicab_%d.txt", num))
        defer vertices_file.Close()

        vertices, _ := io.ReadAll(vertices_file)
        points := []point{}

        str := strings.Split(string(vertices), "\r\n")

        for _, v := range str {
                nums := strings.Split(v, string([]byte{0x09}))
                if len(nums) < 2 {
                        continue
                }
                x, _ := strconv.Atoi(nums[0])
                y, _ := strconv.Atoi(nums[1])
                points = append(points, point{x: uint16(x), y: uint16(y)})
        }

        for i := 0; i < num; i++ {
                incidence = append(incidence, make([]uint16, num))
                distance = append(distance, make([]uint16, num))
        }

        for i, v1 := range points {
                for j, v2 := range points {
                        distance[i][j] = taxi(v1, v2)
                }
        }

        for i := 0; i < num; i++ {
                for j := i + 1; j < num; j++ {
                        edges = append(edges, edge{v1: i, v2: j})
                }
        }

        resChan = make(chan *result, nt)
        for i := 0; i < nt; i++ {
                go processGraph(num, i)
        }
        cost := 0
        for {
                curRes := <-resChan
                if cost < curRes.cost {
                        cost = curRes.cost
                        cnt := make(map[int]int)
                        ecnt := 0
                        for i := 0; i < num; i++ {
                                for j := i + 1; j < num; j++ {
                                        if curRes.incidence[i][j] {
                                                cnt[i] = 1
                                                cnt[j] = 1
                                                ecnt++
                                                fmt.Printf("e %d %d\n", i+1, j+1)
                                        }
                                }
                        }
                        fmt.Printf("c Вес подграфа = %d\n", cost)
                        fmt.Printf("p edge %d %d\n", len(cnt), ecnt)
                }
        }
}
