# Obtiene el último tag
latest_tag=$(git describe --tags --abbrev=0)

# Incrementa el número de revisión
latest_tag=$(git describe --tags --abbrev=0)

# Si el tag es del tipo "pre-release" (por ejemplo, 1.0.2-beta)
if [[ $latest_tag =~ - ]]; then
    # Divide el número de versión y el sufijo
    version=$(echo $latest_tag | awk -F- '{print $1}')
    suffix=$(echo $latest_tag | awk -F- '{print $2}')

    # Incrementa el número de versión
    new_version=$(echo $version | awk -F. '{$NF = $NF + 1;} 1' OFS=.)

    # Combina el nuevo número de versión con el sufijo
    new_tag="${new_version}-${suffix}"

# Si el tag es del tipo "sufijo con punto" (por ejemplo, 1.0.2.test)
elif [[ $latest_tag =~ \. ]]; then
    # Divide el número de versión y el sufijo
    base_version=$(echo $latest_tag | rev | cut -d. -f2- | rev)
    suffix=$(echo $latest_tag | rev | cut -d. -f1 | rev)

    # Verifica si el sufijo es numérico
    if [[ $suffix =~ ^[0-9]+$ ]]; then
        # Incrementa el número de versión
        new_tag=$(echo $latest_tag | awk -F. '{$NF = $NF + 1;} 1' OFS=.)
    else
        # Mantiene el mismo sufijo
        new_tag="${base_version}.${suffix}"
    fi

else
    # Incrementa el número de revisión si no hay sufijo
    new_tag=$(echo $latest_tag | awk -F. '{$NF = $NF + 1;} 1' OFS=.)
fi


docker build . -t sidecar/filter/controller:$new_tag -f deploy/docker/controller/Dockerfile

docker tag sidecar/filter/controller:$new_tag routerdi1315.uv.es:33443/sidecar/filter/controller:$new_tag

docker push routerdi1315.uv.es:33443/sidecar/filter/controller:$new_tag

docker build . -t sidecar/filter/admission:$new_tag -f deploy/docker/admission/Dockerfile

docker tag sidecar/filter/admission:$new_tag routerdi1315.uv.es:33443/sidecar/filter/admission:$new_tag

docker push routerdi1315.uv.es:33443/sidecar/filter/admission:$new_tag