from dataclasses import dataclass

from creature import Creature


@dataclass
class Trait:
    name: str
    description: str
    material: str
    creature: Creature
