import matplotlib.pyplot as plt
import numpy as np


def getMean(nombre_archivo, multiplicador):
    with open(nombre_archivo, "r") as file:
        datos = file.readlines()
        datos = [float(dato.strip()) / 1000 for dato in datos]

        # Reshape los datos para tener una matriz de (50, multiplicador)
        datos = np.array(datos).reshape(50, multiplicador)

        # Calcular la media a lo largo del eje 1 (eje de las columnas)
        medias = np.mean(datos, axis=1)

    return medias


dir = "../../data/coolstart/coolstart-go-8/"

# Obtiene la media para cada archivo
medias_0 = getMean(dir + "0.txt", 8)
medias_1 = getMean(dir + "1.txt", 8)
medias_2 = getMean(dir + "2.txt", 8)
medias_5 = getMean(dir + "5.txt", 8)
medias_10 = getMean(dir + "10.txt", 8)

datos = [medias_0, medias_1, medias_2, medias_5, medias_10]

plt.boxplot(datos)
plt.xticks(
    [1, 2, 3, 4, 5], ["0 Sidecar", "1 Sidecar", "2 Sidecar", "5 Sidecar", "10 Sidecar"]
)
plt.ylabel("Cold Start time (s)")
plt.show()
