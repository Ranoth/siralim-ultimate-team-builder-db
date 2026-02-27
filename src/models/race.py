from dataclasses import dataclass


from creature_class import CreatureClass
from creature import Creature


@dataclass
class Race:
    name: str
    icon: str
    creature_class: CreatureClass
    creatures: list[Creature]