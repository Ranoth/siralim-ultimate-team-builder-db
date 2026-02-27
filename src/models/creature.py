from dataclasses import dataclass

from creature_class import CreatureClass
from trait import Trait


@dataclass
class Creature:
    creature_class: CreatureClass
    race: str
    name: str
    trait: Trait
