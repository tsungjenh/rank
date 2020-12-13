## About The Project

It's a common scenario to order a list of objects in daily programming tasks.  If we want to re-order a list of objects. The intuitive idea is to update the order for the whole array. Something like this:

```
(id, order)
[(1, 1), (2, 2), (3, 3)]  ->   [(4,1), (1,2), (2, 3), (3, 4)]
```

This will take `O(n)` time complexity to re-order which is not acceptable in many applications.
This project intends to solve this problem. Allow insertion between 2 strings with `O(1)` complexity.

## Installation
~~~sh
go get -u github.com/tsungjenh/rank
~~~

## Usage
~~~go
import (
	"github.com/tsungjenh/rank"
	"fmt"
)

func main() {
	r, _ := NewRank(10)

	ranks := r.NewRanks(7)
	fmt.Println(ranks)  // []string{"4DI", "8R", "D4I", "HI", "LVI", "Q9", "UMI"}

	ranksBetween2Nums := r.NewRanksBetween("0", "8", 7)
	fmt.Println(ranksBetween2Nums) // []string{"1", "2", "3", "4", "5", "6", "7"}

	rankBetween := r.Insert("0", "Z")
	fmt.Println(rankBetween	// "HI"

	rankNext := r.Next("A1")
	fmt.Println(rankNext) // "A100000001"

	rankPrev := r.Prev("CAC")
	fmt.Println(rankPrev) // "CABZZZZZZZ"
}
~~~
