package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Point struct {
	x float64
	y float64
}

func calcDistance(a, b Point) float64 {
	return math.Sqrt(math.Pow(a.x-b.x, 2) + math.Pow(a.y-b.y, 2))
}

func initCenters(data []Point, k int) []Point {
	rand.Seed(time.Now().UnixNano())
	centers := make([]Point, 0, k)
	for len(centers) < k {
		pointIdx := rand.Intn(len(data))
		centers = append(centers, data[pointIdx])
	}
	return centers
}

func assignDocs(data []Point, centers []Point) [][]int {
	clusters := make([][]int, len(centers))
	for i, point := range data {
		closestCenterIdx := -1
		closestDistance := math.MaxFloat64
		for j, center := range centers {
			distance := calcDistance(point, center)
			if distance < closestDistance {
				closestDistance = distance
				closestCenterIdx = j
			}
		}
		clusters[closestCenterIdx] = append(clusters[closestCenterIdx], i)
	}
	return clusters
}

func updateCenters(data []Point, clusters [][]int, oldCenters []Point) []Point {
	newCenters := make([]Point, 0, len(oldCenters))
	for i, cluster := range clusters {
		var newCenter Point
		if len(cluster) == 0 {
			// 空のクラスタに対しては既存の代表点を維持
			newCenters = append(newCenters, oldCenters[i])
			continue
		}
		for _, pointIdx := range cluster {
			newCenter.x += data[pointIdx].x
			newCenter.y += data[pointIdx].y
		}
		newCenter.x /= float64(len(cluster))
		newCenter.y /= float64(len(cluster))
		newCenters = append(newCenters, newCenter)
	}
	return newCenters
}

func calcIntraDist(data []Point, centers []Point, clusters [][]int) float64 {
	totalDist := 0.0
	totalPoints := 0
	for i, cluster := range clusters {
		for _, pointIdx := range cluster {
			totalDist += calcDistance(centers[i], data[pointIdx])
			totalPoints++
		}
	}
	return totalDist / float64(totalPoints)
}

func calcInterDist(centers []Point) float64 {
	totalDist := 0.0
	totalPairs := 0
	for i := 0; i < len(centers); i++ {
		for j := i + 1; j < len(centers); j++ {
			totalDist += calcDistance(centers[i], centers[j])
			totalPairs++
		}
	}
	return totalDist / float64(totalPairs)
}

func main() {
	file, err := os.Open("data1.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	var data []Point
	var names []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), "\t")
		if len(parts) == 3 {
			names = append(names, parts[0])
			latitude, _ := strconv.ParseFloat(parts[1], 64)
			longitude, _ := strconv.ParseFloat(parts[2], 64)
			data = append(data, Point{latitude, longitude})
		}
	}

	k := 8
	centers := initCenters(data, k)

	for {
		clusters := assignDocs(data, centers)
		newCenters := updateCenters(data, clusters, centers)

		if fmt.Sprintf("%v", newCenters) == fmt.Sprintf("%v", centers) {
			break
		}

		centers = newCenters
	}

	// Step 1: 代表点の初期化
	fmt.Println("Step 1. 代表点の初期化")
	centers = initCenters(data, k)
	for i, center := range centers {
		fmt.Printf("クラスタ%dの代表点: [%f, %f]\n", i+1, center.x, center.y)
	}

	var steps int
	for {
		// Step 2: クラスタ割り当て
		fmt.Println("Step 2. クラスタ割り当て")
		clusters := assignDocs(data, centers)
		for i, cluster := range clusters {
			fmt.Printf("Cluster %d: ", i+1)
			for _, pointIdx := range cluster {
				fmt.Printf("%s, ", names[pointIdx])
			}
			fmt.Println()
		}

		// Step 3: 代表点の更新
		fmt.Println("Step 3. 代表点の更新")
		newCenters := updateCenters(data, clusters, centers)
		for i, center := range newCenters {
			fmt.Printf("クラスタ%dの代表点: [%f, %f]\n", i+1, center.x, center.y)
		}

		// 代表点が変わらないかチェック
		if fmt.Sprintf("%v", newCenters) == fmt.Sprintf("%v", centers) {
			fmt.Println("代表点が変化しなかったので処理を終了")
			break
		}

		centers = newCenters
		steps++
	}

	// クラスタリング結果評価
	fmt.Println("クラスタリング結果評価")
	Sintra := calcIntraDist(data, centers, assignDocs(data, centers))
	Sinter := calcInterDist(centers)
	fmt.Printf("クラスタ内分散:%f\n", Sintra)
	fmt.Printf("クラスタ間分散:%f\n", Sinter)
	fmt.Printf("クラスタリング結果の評価値:%f\n", Sinter/Sintra)
}
