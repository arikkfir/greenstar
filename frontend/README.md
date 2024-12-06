# GreenSTAR Frontend

This is the GreenSTAR frontend application served to users. It is packaged into a container and runs as a stateless set
of pods running:

* Production: NGINX serving statically-built React SPA
* Workstation: Node.js running Vite Server serving SPA directly with live reload

This setup supports horizontal scale-out if needed.

## Architecture & Structure

The frontend application follows a structured architecture aimed at modularity and maintainability. Here's a brief
overview:

### Architecture

- **Single Page Application (SPA)**: Built using Vite & React, this application provides a dynamic and seamless user
  experience by serving content dynamically without the need for page reloads.
- **State Management**: Utilizes the Context & Component API for managing the application state.
- **Routing**: React Router is used for client-side routing, allowing navigation among views of various components in a
  single-page application.
- **UI Components**: Built with Material-UI (MUI) for reusable, sustainable, and customizable components.

### Structure

The important parts of the application in the file-system are the following:

- **`Dockerfile.frontend`**: Packages the application into a container image. Contains the following targets:
    - **`development`**: Used when developing locally, and results in a container that runs the Vite Server directly
      serving the application and utilizing HMR for dynamic reloading on changes (integrated with Skaffold.)
    - **`production`**: Used by the CI workflows to generate a container that converts the application to a set of
      static assets that are served by an NGINX server (with support for HTML5 routing.)
- **`nginx.conf`**: Configuration for the NGINX server, when generating production-oriented container image.
- **`src/`**: The main source directory for application code:
  - **`main.tsx`**: Entry-point of the application.
  - **`App.tsx`**: The root component of the application.
  - **`client/`**: Contains generated REST client code as well as non-generated related code.
  - **`components/`**: Contains reusable React components.
  - **`hooks/`**: Contains custom React hooks for reusing logic.
  - **`pages/`**: Contains page components that represent different routes of the application.
  - **`providers/`**: Contains reusable React context providers.
  - **`util/`**: Contains miscellaneous utilities.

### Development Environment

- **Build Tools**: Vite is used for development to speed up the build times with live reloads and efficient module
  replacement.
- **Production**: In production, the application is statically built and served using NGINX for optimized performance.

This architecture is designed for scalability and ease of maintenance, allowing teams to work on different modules
concurrently and seamlessly integrate changes.

Please see the [CONTRIBUTING](../CONTRIBUTING.md) guide for more information.

## General open items

- [ ] Generate client hooks into the new `hooks` folder
- [ ] Switch to file-based routes (see [File Routes](https://reactrouter.com/how-to/file-route-conventions))
