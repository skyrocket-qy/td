using UnityEngine;

[CreateAssetMenu(fileName = "NewHeroClass", menuName = "Game/Progression/Hero Class")]
public class HeroClassData : ScriptableObject
{
    public string className = "New Class";
    // Assign the PassiveTreeData asset that defines the passive tree for this hero class.
    public PassiveTreeData passiveTree; 

    // You can add other class-specific properties here, like base stats, starting skills, etc.
}
