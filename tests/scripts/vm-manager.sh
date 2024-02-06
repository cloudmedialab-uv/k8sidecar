#!/bin/bash

# REMOTE_USER, REMOTE_HOST, PRIVATE_KEY_PATH shoud be env vars sourced by vars file

start_vm() {
    for vm in "${@:1}"; do
        echo "Intentando encender $vm..."

        # Comprobar si la máquina virtual ya está en ejecución
        vm_status=$(ssh -i "$REMOTE_PRIVATE_KEY" "$REMOTE_USER"@"$REMOTE_HOST" 'bash -l -c "virsh list --name | grep '$vm'"')
        if [ ! -z "$vm_status" ]; then
            echo "$vm ya está activa. Saltando..."
            continue
        fi

        success=0
        while [ $success -eq 0 ]; do
            ssh -i "$REMOTE_PRIVATE_KEY" "$REMOTE_USER"@"$REMOTE_HOST" 'bash -l -c "virsh start '$vm'"'
            if [ $? -eq 0 ]; then
                echo "$vm ha sido encendida con éxito."
                success=1
            else
                echo "Fallo al encender $vm. Reintentando en 5 segundos..."
                sleep 5
            fi
        done
    done
}



stop_vm() {
    for vm in "${@:1}"; do
        echo "Apagando $vm..."
        ssh -i "$REMOTE_PRIVATE_KEY" "$REMOTE_USER"@"$REMOTE_HOST" 'bash -l -c "virsh shutdown '$vm'"'
    done
}

if [ $# -lt 2 ]; then
    echo "Uso: $0 {start|stop} vm1 [vm2 vm3 ...]"
    exit 1
fi

up_vm() {
    local desired_vms=("$@")

    # Obtiene una lista de todas las máquinas virtuales en ejecución
    local running_vms=($(ssh -i "$REMOTE_PRIVATE_KEY" "$REMOTE_USER"@"$REMOTE_HOST" 'bash -l -c "virsh list --name"'))

    # Apaga las VMs que no están en la lista deseada
    for running_vm in "${running_vms[@]}"; do
        if ! printf '%s\n' "${desired_vms[@]}" | grep -q -P "^$running_vm$"; then
            echo "Apagando $running_vm..."
            ssh -i "$REMOTE_PRIVATE_KEY" "$REMOTE_USER"@"$REMOTE_HOST" "bash -l -c 'virsh shutdown $running_vm'"

            # Espera y comprueba que la VM se haya apagado
            shut_down=0
            while [ $shut_down -eq 0 ]; do
                sleep 5
                is_running=$(ssh -i "$REMOTE_PRIVATE_KEY" "$REMOTE_USER"@"$REMOTE_HOST" "bash -l -c 'virsh list --name | grep $running_vm'")
                if [ -z "$is_running" ]; then
                    echo "$running_vm ha sido apagada con éxito."
                    shut_down=1
                else
                    echo "Esperando a que $running_vm se apague..."
                fi
            done
        fi
    done

    # Enciende las VMs deseadas que no están en ejecución
    for desired_vm in "${desired_vms[@]}"; do
        if ! printf '%s\n' "${running_vms[@]}" | grep -q -P "^$desired_vm$"; then
            echo "Encendiendo $desired_vm..."
            ssh -i "$REMOTE_PRIVATE_KEY" "$REMOTE_USER"@"$REMOTE_HOST" "bash -l -c 'virsh start $desired_vm'"

            # Espera y verifica que la VM se haya encendido
            started_up=0
            while [ $started_up -eq 0 ]; do
                sleep 5
                is_running=$(ssh -i "$REMOTE_PRIVATE_KEY" "$REMOTE_USER"@"$REMOTE_HOST" "bash -l -c 'virsh list --name | grep $desired_vm'")
                if [ ! -z "$is_running" ]; then
                    echo "$desired_vm ha sido encendida con éxito."
                    started_up=1
                else
                    echo "Esperando a que $desired_vm se encienda..."
                fi
            done
        else
            echo "$desired_vm ya está en ejecución."
        fi
    done
}

ACTION="$1"
shift 

case "$ACTION" in
    start)
        start_vm "$@"
        ;;
    stop)
        stop_vm "$@"
        ;;
    up)
        up_vm "$@"
        ;;
    *)
        echo "Acción no reconocida. Uso: $0 {start|stop|up} vm1 [vm2 vm3 ...]"
        exit 1
        ;;
esac
