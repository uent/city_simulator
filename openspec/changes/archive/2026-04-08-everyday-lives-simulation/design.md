## Context

El motor de simulación soporta cualquier escenario que provea `characters.yaml`, `world.yaml`, y `scenario.yaml` en el formato correcto. El scenario loader ya lee estos archivos y los sirve al engine sin cambios. Este cambio es puro contenido: ningún archivo Go se modifica.

Las simulaciones existentes (default, honey-heist, doom-hell-crusade) establecen el patrón de contenido: mundo con concepto + reglas, personajes con arquitectura psicológica completa, y un scenario.yaml mínimo para overrides de runtime.

La diferencia de esta simulación es de naturaleza, no de formato:
- Sin género externo (no hay crimen, robo, ni amenaza demoníaca)
- Sin game_director — la tensión emerge de las personas, no de un narrador
- Sin misión colectiva — cada personaje tiene su propia vida en curso

## Goals / Non-Goals

**Goals:**
- Crear un escenario de vida cotidiana con 5 personajes profundamente humanos en un edificio de apartamentos
- Demostrar que el motor genera interacciones significativas sin tensión externa ni director de juego
- Establecer un nuevo patrón de escenario: drama relacional sin género

**Non-Goals:**
- No se modifica ningún componente del motor
- No se introduce ninguna mecánica nueva al engine o al loader
- No es un mundo fantástico, alternativo, ni de época — es contemporáneo y mundano por diseño

## Decisions

### Mundo: edificio de apartamentos + barrio inmediato

El espacio compartido forzado es el mecanismo que hace posible las interacciones entre personas que no se eligieron entre sí. Un edificio genera roce inevitables (ascensor, lobby, azotea) sin requerir ninguna justificación narrativa.

**Alternativa considerada**: Plaza de barrio como punto de encuentro. Descartada porque los personajes podrían nunca cruzarse orgánicamente — el edificio garantiza al menos tres puntos de contacto inevitables.

### Personajes: 5 residentes jóvenes (20–30 años) en momentos de vida no resueltos

Cinco es el número que permite tensiones diádicas variadas sin que el elenco sea difícil de seguir. Todos tienen entre 20 y 30 años — una franja donde las personas están construyendo su identidad adulta sin haberla terminado, lo que los hace especialmente permeables a los pequeños eventos de un domingo. Los momentos de tránsito propios de esta edad (primer trabajo serio, primera relación que dura, primera pérdida real, primera decepción con uno mismo) generan fricción sin necesitar backstories elaborados.

**Principio de diseño**: Los personajes no deben ser arquetipos funcionales sino personas cuya complejidad produce comportamiento impredecible. Un personaje que "necesita conexión pero la evita" genera más interacciones interesantes que uno que simplemente "es introvertido".

### Sin game_director

Los escenarios existentes con mayor tensión narrativa usan un director. Este escenario lo omite deliberadamente para explorar si la arquitectura psicológica de los personajes es suficiente para generar drama sin andamiaje externo.

**Riesgo aceptado**: Sin director, la simulación puede volverse episódica o sin arco claro. Esto es aceptable — la vida cotidiana real no tiene arco, y el objetivo es precisamente eso.

### Nombre del directorio: `vida-cotidiana`

Descriptivo, en español (consistente con el estilo del proyecto), y sin ambigüedad sobre el tipo de escenario.

## Risks / Trade-offs

- **Sin tensión externa explícita, las interacciones pueden volverse triviales** → Mitigación: Los personajes tienen tensiones internas no resueltas (separación no procesada, relación que se deshace, deuda financiera, etc.) que crean fricción incluso en conversaciones ordinarias.

- **Sin director, no hay mecanismo para introducir complicaciones** → Mitigación: Los initial_events plantan semillas que los personajes procesan de forma distinta (corte de luz, una nota en el buzón, el mercado del domingo).

- **Los personajes pueden no tener razones para interactuar** → Mitigación: Los espacios compartidos del mundo (lobby, ascensor, azotea) funcionan como catalizadores inevitables; el mundo.yaml los establece como puntos de cruce naturales.
