import json
import random

'''
Data generated as follows:
    K sets of Ni UIDs for K districts with populations of size Ni

    for each uid: pick 1-2 random times in 0,36
        possibilities
            safe
            distress -> safe
            deceased
            distress -> deceased

        times should be drawn from logistic distribution that starts at 0 ends at 36

    sort events by time, save to text file
'''

class District:
    def __init__(self, name, population, center):
        self.name = name
        self.center = center
        self.count = population

def makeDistricts(districts):
    return {dname: District(dname, districts[dname]['pop'], districts[dname]['center']) for dname in districts}

class Event:
    def __init__(self, uid, state, time, loc):
        self.uid = uid
        self.state = state
        self.time = time
        self.locx = loc[0]
        self.locy = loc[1]
    
    def __lt__(self, other):
        return self.time < other.time
    
    def __repr__(self):
        return '{} {} {} {} {}'.format(self.uid, self.state, self.time, self.locx, self.locy)

class User:
    finalState = ['safe'] * 95 + ['deceased'] * 5
    numStates = [2] * 40 + [1] * 60
    maxTime = 36

    def __init__(self, uid, loc):
        self.uid = uid
        self.loc = loc

    def generate(self):
        numEvents = random.choice(User.numStates)

        if numEvents == 1:
            time = random.randint(0, User.maxTime)
            state = random.choice(User.finalState)
            return [Event(self.uid, state, time, self.loc)]
        
        time = random.randint(0, User.maxTime-1)
        first = Event(self.uid, 'distress', time, self.loc)

        time = random.randint(time, User.maxTime)
        secondState = random.choice(User.finalState)
        second = Event(self.uid, secondState, time, self.loc)

        return [first, second]


class DataGenerator:
    def __init__(self, districts):
        currUID = 100
        self.users = []
        for dname in districts:
            for i in range(districts[dname].count):
                u = User(currUID, districts[dname].center)
                currUID += 1
                self.users.append(u)
    
    def generate(self):
        events = []
        for user in self.users:
            for e in user.generate():
                events.append(e)

        events.sort()

        with open('eventsGenerated.txt', 'w') as f:
            for event in events:
                f.write('{}\n'.format(event))


if __name__ == '__main__':
    with open('districts.json') as jf:
        districts = makeDistricts(json.load(jf))
    
    generationMachine = DataGenerator(districts)
    generationMachine.generate()