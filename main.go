package main

import (
	"encoding/json"
	"encoding/binary"
	"os"
	"fmt"
	"net/http"
	"strconv"
	"math"
)

var qTable [19683]float32

func main() {
	LoadQTable()
	http.HandleFunc("/api/hello", helloHandler)
	http.HandleFunc("/api/state_action", stateActionHandler)
	fmt.Println("Server listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	// CORSヘッダーを設定
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	
	// OPTIONSリクエストに対してはヘッダーのみを返して終了
	if r.Method == http.MethodOptions {
		return
	}

	// レスポンスを設定
	response := map[string]string{"message": "Hello, World!"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func stateActionHandler(w http.ResponseWriter, r *http.Request) {
	// CORSヘッダーを設定
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	
	// OPTIONSリクエストに対してはヘッダーのみを返して終了
	if r.Method == http.MethodOptions {
		return
	}

	// リクエストボディをパースして数値を取得
	stateIndexStr := r.URL.Query().Get("state_index")
	if stateIndexStr == "" {
		http.Error(w, "Missing 'state' parameter", http.StatusBadRequest)
		return
	}
	stateIndex, err := strconv.Atoi(stateIndexStr)
	if err != nil {
		http.Error(w, "Invalid 'state' parameter", http.StatusBadRequest)
		return
	}
	fmt.Println(stateIndex)
	state := Index2State(stateIndex)
	PrintState(state)

	optimumAction := CalculateOptimumAction(state, CalculateAvailableAction(state))
	fmt.Println(optimumAction)

	action := optimumAction[1]*3 + optimumAction[0]

	// レスポンスを設定
	response := map[string]int{"action": action}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func State2Index(state [3][3]rune) int {
	index := 0
	for i := 0; i < 9; i++ {
		var value int
		cell := state[i/3][i%3]

		switch cell {
		case '_':
			value = 0
		case 'o':
			value = 1
		case 'x':
			value = 2
		default:
			fmt.Printf("State2Index() state[%d][%d] is invalid '%c'\n", i/3, i%3, cell)
		}

		// 3 のべき乗計算 (3^i)
		index += value * int(math.Pow(3, float64(i)))
	}
	return index
}

func Index2State(state_index int) [3][3]rune {
	state := [3][3]rune{
		{'_', '_', '_'},
		{'_', '_', '_'},
		{'_', '_', '_'},
	}

	for i := 0; i < 9; i++ {
		switch state_index % 3 {
		case 0:
			state[i/3][i%3] = '_'
		case 1:
			state[i/3][i%3] = 'o'
		case 2:
			state[i/3][i%3] = 'x'
		}
		state_index /= 3
	}

	return state
}

func PrintState(state [3][3]rune) {
	fmt.Printf("%c %c %c\n", state[0][0], state[1][0], state[2][0])
	fmt.Printf("%c %c %c\n", state[0][1], state[1][1], state[2][1])
	fmt.Printf("%c %c %c\n", state[0][2], state[1][2], state[2][2])
}

func CalculateAvailableAction(state [3][3]rune) [][2]int {
	var availableCells [][2]int
	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			if state[x][y] == '_' {
				availableCells = append(availableCells, [2]int{x, y})
			}
		}
	}
	return availableCells
}

func CalculateAfterState(state [3][3]rune, action [2]int) [3][3]rune {
	state[action[0]][action[1]] = 'x'
	return state
}

func CalculateOptimumAction(currentState [3][3]rune, availableActionCells [][2]int) [2]int {
	if len(availableActionCells) == 1 {
		// 唯一のアクションを返す
		return availableActionCells[0]
	}

	// 初期化: 最初のアクションを最大アクションとする
	maxAction := availableActionCells[0]

	// 残りのアクションを比較し、最適なものを選ぶ
	for i := 1; i < len(availableActionCells); i++ {
		if qTable[State2Index(CalculateAfterState(currentState, availableActionCells[i]))] > qTable[State2Index(CalculateAfterState(currentState, maxAction))] {
			maxAction = availableActionCells[i]
		}
	}

	return maxAction
}

func LoadQTable() {
	// バイナリファイルの読み込み
	file, err := os.Open("q_table_go")
	if err != nil {
		panic(err)
	}

	// バイナリファイルを読み込んでfloat32のスライスに変換
	err = binary.Read(file, binary.LittleEndian, &qTable)
	if err != nil {
		panic(err)
	}
	defer file.Close()
}