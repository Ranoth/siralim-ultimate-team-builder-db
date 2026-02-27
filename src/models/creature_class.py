from dataclasses import dataclass

from creature import Creature

@dataclass
class CreatureClass:
    name: str
    icon: str
    creatures: list[Creature]