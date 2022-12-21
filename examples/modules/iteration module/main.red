IMPORT iter.mred i

// defining some values for the array
PUSH 1
PUSH 2
PUSH 3
PUSH 9
PUSH 5

// making and storing array
MAKEARRAY
MODSTORE i array

// intialise iteration
MODRUN i init

// run function
MODRUN i iter isLooping

// reset module
MODRUN i reset