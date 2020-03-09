from block import Block
from time import time

def main():
    chain = []
    block = Block(len(chain), time(), 'a', 0)
    chain.append(block)
    block = Block(len(chain), time(), 'a', chain[-1].hash)
    chain.append(block)
    print(chain)
    


if __name__ == '__main__':
    main()