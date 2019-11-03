import numpy as np
import cv2
import matplotlib as mpl
mpl.use('TkAgg')
mpl.rcParams['figure.dpi']= 300
import matplotlib.pyplot as plt
import cv2
import json



'''
Troy population- 49374
find centers of districts with distance to border 
    or pick n centers
assign population to districts
    pick #
    for each person in a district, generate their distance from center with a Gaussian

distribution
49374
downtown        8174
the hill        8000
south central   7500
south Troy      7000
Eastside        6700 
Sycaway         6000
Lansingburgh    4000
Frear park      2000
'''

class Event:
    def __init__(self, uid, status, time, xloc, yloc):
        self.uid = int(uid)
        self.status = status
        self.time = int(time)
        self.loc = np.array([float(xloc), float(yloc)])

class City:
    def __init__(self, im, districts):
        self.im = im
        self.districts = districts
    
    def plotDistricts(self, time = 0):
        plt.axis('off')
        plt.imshow(im)

        for dname in self.districts:
            self.districts[dname].plotPop()
        
        #should clear plt
        #plt.show()
        plt.savefig('Troy_{}_hours.png'.format(time))
        plt.clf()

    def processEvent(self, event):
        dname = self.findClosestDistrict(event.loc)
        districts[dname].processEvent(event)

        return event.time


    #Neareast Neighbor, in context of image, not lat/lng
    def findClosestDistrict(self, point):
        names = list(self.districts.keys())
        dists = []

        for dname in names:
            dists.append( np.sum( (self.districts[dname].center - point) ** 2) **0.5 )
        
        dists = np.array(dists)
        min_i = np.argmin(dists)
    
        return names[min_i]

class District:
    def __init__(self, name, population, center, covar):
        self.name = name
        self.center = center
        self.covar = covar
        self.count = population
        self.unknownCount = self.count
        self.distressUID = set()
        self.safeUID = set()
        self.deceasedUID = set()
    
    def plotPop(self):
        popUNK = np.random.multivariate_normal(self.center, self.covar, size = (self.unknownCount,))
        popHELP = np.random.multivariate_normal(self.center, self.covar, size = (len(self.distressUID),))
        popSAFE = np.random.multivariate_normal(self.center, self.covar, size = (len(self.safeUID),))
        popDEAD = np.random.multivariate_normal(self.center, self.covar, size = (len(self.deceasedUID),))

        plt.scatter(popUNK[:,0], popUNK[:,1], c = "GREEN", s = 0.1)
        plt.scatter(popDEAD[:,0], popDEAD[:,1], c = "BLACK", s = 0.1)
        plt.scatter(popSAFE[:,0], popSAFE[:,1], c = "BLUE", s = 0.1)
        plt.scatter(popHELP[:,0], popHELP[:,1], c = "RED", s = 0.1)
    
    def processEvent(self, event):
        if event.uid not in self.distressUID and event.uid not in self.safeUID and event.uid not in self.deceasedUID:
            self.unknownCount -= 1

        if event.status == 'deceased':
            self.deceasedUID.add(event.uid)
            if event.uid in self.safeUID:
                self.safeUID.remove(event.uid)

            if event.uid in self.distressUID:
                self.distressUID.remove(event.uid)

        elif event.status == 'safe':
            self.safeUID.add(event.uid)

            if event.uid in self.distressUID:
                self.distressUID.remove(event.uid)

        elif event.status == 'distress':
            self.distressUID.add(event.uid)

            if event.uid in self.safeUID:
                self.safeUID.remove(event.uid)

        else:
            print('ERROR: status {} not recognized, must be one of deceased, safe, distress'.format(event.status))


        if self.unknownCount == 0:
            pass
            #print('District {} is finished!'.format(self.name))
        
        if self.unknownCount < 0:
            print('ERROR: Someone switched districts!!')


def makeDistricts(districts):
    return {dname: District(dname, districts[dname]['pop'], districts[dname]['center'], districts[dname]['covar']) for dname in districts}


if __name__ == '__main__':
    im = cv2.imread('troy.png', cv2.IMREAD_GRAYSCALE)
    im = np.stack((im,)*3, axis=-1)
    #plt.imshow(im)

    with open('districts.json') as jf:
        districts = makeDistricts(json.load(jf))

    Troy = City(im, districts)

    time = 0

    with open('eventsGenerated.txt') as txt:
        events = list(  map(lambda line: Event(*line.split()), txt.readlines())  )

    Troy.plotDistricts()

    lastUpdate = 0
    updateInterval = 6
    #assume events are sorted
    for event in events:
        time = Troy.processEvent(event)
        if time // updateInterval > lastUpdate:
            Troy.plotDistricts(time = time)
            lastUpdate = time // 6
    
    Troy.plotDistricts(time = 36)