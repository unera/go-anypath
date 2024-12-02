// anypath is an navigator through any object

package anypath

import (
	"fmt"
	"strconv"
)

func parsePath(path string) ([]any, error) {

	data := path + string(rune(0))
	cs, p, pe := 0, 0, len(data)
	eof := pe 
	ts, te := 0, 0
	found := make([]any, 16) // 16 is enough for ANY cases
	foundIndex := -1

	pushFound := func (e any) {
		foundIndex++

		if foundIndex >= len(found) {
			nfound := make([]any, len(found) + 16)
			copy(nfound, found)
			found = nfound
		}

		found[foundIndex] = e
	}

	%%{
		machine anypath;
		write data;
		
		eop = 0;

		index = ("-"? digit+)
			>{ ts, te = p, 0 }
			%{ te = p };

		divider = '.';
		token = (any - space - [\[\]\.] - eop)+
			>{ ts, te = p, 0 }
			%{
				te = p
			   	pushFound(string(data[ts:te]))
 			};

		bracketIndex = ('[' index ']')
			%{{
				index, err := strconv.ParseInt(string(data[ts:te]), 10, 64)
				if err != nil {
					return nil, err
				}
				pushFound(index)
			}}
		;



		path = (token? (bracketIndex | (divider  token))* eop)
			$err{
				if p == eof - 1 {
					return nil, fmt.Errorf("Unexpected EOF at pos=%d", p+1)
				}

				return nil, fmt.Errorf("Unexpected symbol '%c' at pos %d", data[p], p+1)
			}

		;

		anypath := path;
		write init;
		write exec;
	}%%


	return found[:foundIndex + 1], nil
}


