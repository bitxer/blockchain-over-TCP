class Block:
    def __init__(self, index, timestamp, data, parentHash):
        self.index = index
        self.timestamp = timestamp
        self.data = data
        self.parentHash = parentHash
        self.genHash()
    
    def __repr__(self):
        return '<Block: index={} timestamp={} data={} parentHash={} hash={}>'.format(self.index, self.timestamp, self.data, self.parentHash, self.hash)

    def genHash(self):
        self.hash = hash('{}{}{}'.format(self.index, self.timestamp, self.data, self.parentHash))
