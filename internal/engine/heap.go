package engine

type OrderPriorityQueue []*Order

type BuyPriorityQueue struct {
		OrderPriorityQueue
}

type SellPriorityQueue struct {
		OrderPriorityQueue
}

func (pq OrderPriorityQueue) Len() int {
		return len(pq)
}

func (pq BuyPriorityQueue) Less(i, j int) bool {
		if pq.OrderPriorityQueue[i].Price == pq.OrderPriorityQueue[j].Price {
				return pq.OrderPriorityQueue[i].In < pq.OrderPriorityQueue[j].In
		}
		return pq.OrderPriorityQueue[i].Price > pq.OrderPriorityQueue[j].Price
}

func (pq SellPriorityQueue) Less(i, j int) bool {
		if pq.OrderPriorityQueue[i].Price == pq.OrderPriorityQueue[j].Price {
				return pq.OrderPriorityQueue[i].In < pq.OrderPriorityQueue[j].In
		}
		return pq.OrderPriorityQueue[i].Price < pq.OrderPriorityQueue[j].Price
}

func (pq OrderPriorityQueue) Swap(i , j int){
		pq[i], pq[j] = pq[j], pq[i]
		pq[i].Index=i
		pq[j].Index=j
}

func (pq *OrderPriorityQueue) Push(x interface{}) {
		item := x.(*Order)
		*pq = append(*pq, item)
}

func (pq *OrderPriorityQueue) Pop() interface{} {
		old := *pq
		n := len(*pq)
		item := old[n-1]
		old[n-1]=nil
		*pq = old[0:n-1]
		return item
}

func (pq OrderPriorityQueue) Top() *Order{
		if len(pq) == 0 {
				return nil
		}
		return pq[0]
}
