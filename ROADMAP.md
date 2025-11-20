# Roadmap: CollabSphere (Chat y Pizarra Colaborativa)

Este documento sigue el progreso del desarrollo del proyecto.

---

### Fase 1: El Esqueleto - Estructura y Fundamentos

*   [x] Tarea 1: Estructura del Proyecto
*   [x] Tarea 2: Inicialización del módulo
*   [x] Tarea 3: Elección e instalación de librerías
*   [x] Tarea 4: Servidor HTTP Básico
*   [x] Tarea 5: Endpoint de Health Check
*   [x] Tarea 6: Gestión de Configuración

---

### Fase 2: El Corazón - El Hub de WebSockets

*   [x] Tarea 1: El `Upgrader` de WebSocket
*   [x] Tarea 2: La struct `Client`
*   [x] Tarea 3: La struct `Hub`
*   [x] Tarea 4: Goroutines del `Hub` (`run` method)

---

### Fase 3: La Primera Conversación - Chat Básico

*   [x] Tarea 1: Modelo de Mensaje (struct)
*   [x] Tarea 2: Goroutines del Cliente (`readPump` y `writePump`)
*   [x] Tarea 3: Lógica de Broadcasting en el `Hub`
*   [x] Tarea 4: Notificaciones de Conexión/Desconexión

---

### Fase 4: Múltiples Mundos - Salas de Chat y Mensajes Directos

*   [x] Tarea 1: Gestión de Salas en el `Hub`
*   [x] Tarea 2: API de Salas (unirse a una sala)
*   [x] Tarea 3: Enrutamiento de Mensajes por sala
*   [x] Tarea 4: Mensajes Directos (Bonus)

---

### Fase 5: La Memoria - Persistencia de Datos

*   [x] Tarea 1: Elección e integración de Base de Datos
*   [x] Tarea 2: Capa de Repositorio
*   [x] Tarea 3: Esquema de Base de Datos
*   [x] Tarea 4: Integración con la lógica de negocio

---

### Fase 6: El Plus - Pizarra Colaborativa

*   [ ] Tarea 1: Nuevos Tipos de Mensajes para dibujo
*   [ ] Tarea 2: Lógica de retransmisión de eventos de dibujo
*   [ ] Tarea 3: Sincronización del estado de la pizarra

---

### Fase 7: El Blindaje - Seguridad y Producción

*   [ ] Tarea 1: Autenticación (JWT)
*   [ ] Tarea 2: Autorización
*   [ ] Tarea 3: Logging Estructurado
*   [ ] Tarea 4: Tests Unitarios
*   [ ] Tarea 5: Dockerfile
*   [ ] Tarea 6: Graceful Shutdown
