import matplotlib.pyplot as plt
import json
from statistics import mean

def clearData(data,lim=30):
    return list(filter(lambda x:x<lim,data))

def getLatency(name):
    with open(name, "r") as file:
        times=[]
        for line in file:
            data = json.loads(line)
            latency= data["tsFinalTime"]-data["tsReceivedTime"]
            
            times.append(latency)

    return times


dir = "../../data/latency-envoy/10/"

# Obtiene la media para cada archivo
datos_0 = getLatency("../../data/latency-go/0.json")
datos_1 = getLatency(dir + "1.json")
datos_2 = getLatency(dir + "2.json")
datos_5 = getLatency(dir + "5.json")
datos_10 = getLatency(dir+"10.json")


#datos = [datos_0, datos_10]

print("media de datos_0: ", mean(datos_0))
print("media de datos_1: ", mean(datos_1))
print("media de datos_2: ", mean(datos_2))
print("media de datos_5: ", mean(datos_5))
print("media de datos_10: ", mean(datos_10))

print("diferencia entre 0 y 10", mean(datos_10)-mean(datos_0))



datos = [clearData(datos_0), 
        clearData(datos_1),
        clearData(datos_2), 
        clearData(datos_5), 
        clearData(datos_10)]


plt.boxplot(datos)
plt.xlabel("Number of filters", fontsize=24)
plt.tick_params(axis="x",labelsize=20)
plt.tick_params(axis="y",labelsize=15)

plt.xticks([1, 2, 3, 4, 5], ["0", "1", "2", "5", "10"])
#plt.xticks([1, 2], ["0", "10"])

plt.ylabel("Latency time (ms)", fontsize=24)
plt.tight_layout()
plt.savefig("../../plots/latency-envoy.png")