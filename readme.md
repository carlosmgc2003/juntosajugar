# Ingeniería de Software Avanzado 2020
![JaJ Logo](./img/juntosajugarlogo.png)
## Juntos A Jugar
> Proyecto de aplicación web didáctica para la materia Ing. SW Avanzado

### Setup del ambiente de desarrollo
#### Requisitos
- Docker y Docker Compose

#### Instrucciones
1. ``git clone https://github.com/carlosmgc2003/juntosajugar.git``
2. ``cd juntosajugar``
3. ``docker-compose up``

#### Que despliega
- Un contenedor MySQL.
- Un contenedor Adminer para visualizar el contenido de la BD.
- El contenedor de la API Web, con el código fuente incluido en ``/go/src/juntosajugar``

#### Para probar la api
Hay disponibles ejemplos de Request JSON en la carpeta ``juntosajugar/test/http``
