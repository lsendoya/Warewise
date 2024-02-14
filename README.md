# Warewise

## Descripción
Warewise es un sistema de gestión de  productos. Este proyecto implementa una arquitectura hexagonal, lo que permite una separación clara entre la lógica de negocio y los adaptadores de infraestructura, mejorando así la mantenibilidad y la escalabilidad del sistema.

## Arquitectura y Tecnologías
El sistema Warewise está construido utilizando una variedad de tecnologías y servicios de AWS para proporcionar una solución robusta y escalable:

### Arquitectura Hexagonal
Utilizamos la arquitectura hexagonal para organizar el código de Warewise, lo que nos permite aislar la lógica de negocio de los detalles de la infraestructura y las interfaces de usuario. Esto facilita la prueba de la lógica de negocio de manera aislada y mejora la flexibilidad del sistema.

### AWS Lambda y API Gateway
Para la API, Warewise utiliza AWS Lambda y API Gateway, proporcionando una solución sin servidor altamente escalable. Esto permite que Warewise maneje eficientemente las solicitudes de API con una gestión mínima de la infraestructura.

### Secret Manager
Los secretos de la base de datos y otras credenciales sensibles son gestionados de forma segura utilizando AWS Secret Manager. Esto asegura que la información crítica esté protegida y accesible solo para los servicios autorizados.

### AWS Cognito
Warewise utiliza AWS Cognito para la autenticación y autorización de usuarios, ofreciendo un servicio seguro y escalable para el manejo de identidades de usuario. Esto permite una gestión de acceso y control de identidad eficiente, asegurando que solo los usuarios autorizados puedan acceder a ciertas funcionalidades del sistema.

### Base de Datos PostgreSQL de AWS
Warewise utiliza PostgreSQL de AWS como su sistema de gestión de base de datos. La integración con Secret Manager permite una conexión segura y confiable a la base de datos, facilitando la gestión eficiente de los datos de inventario.

### AWS S3
Para el almacenamiento de fotos de los productos, Warewise aprovecha AWS S3, ofreciendo una solución de almacenamiento escalable y segura. Esto permite a los usuarios cargar y acceder a imágenes de productos fácilmente.


