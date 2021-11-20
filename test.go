package main

import "sync"

type Change struct {
	leguanlock bool
	data interface{}
	isFull bool
	//list []
	l sync.Mutex
}
func (c Change) get( ) interface{}{
	if c.isFull {
		for {
			if c.leguanlock {
				continue

			}else{
				c.isFull = false
				return  c.data
			}
		}
	} else {
		//block;
		//gopark()
		// 等待 isfull 条件满足
		if c.isFull {

		}
	}
}

func ( c Change ) set( in interface{}){
	if c.isFull {
		// block;
	} else if !c.leguanlock {
		c.leguanlock = true
		c.data = in
		c.leguanlock = false
		c.isFull = true
	}

}


select name, project, score from ( select )