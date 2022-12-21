// defining some values for the array
PUSH 1
PUSH 2
PUSH 3
PUSH 9
PUSH 5

// making and storing array
MAKEARRAY
STORE array

// initialise iteration variable
PUSH 0
STORE n

FUNC iter
// load array value at index n and print
LOAD array n
PRINT

// increase n by 1
LOAD n
PUSH 1
ADD
STORE n

// check whether n is equal to the length of the array and end function if it is
LOAD array
LEN
LOAD n
EQ
NOT
STORE isLooping

ENDFUNC

// initialise condition variable
PUSH true
STORE isLooping

// run using condition variable
RUN iter isLooping
