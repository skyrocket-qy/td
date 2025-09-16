using UnityEngine;

[CreateAssetMenu(fileName = "NewAscendancyClass", menuName = "Game/Progression/Ascendancy Class")]
public class AscendancyClassData : ScriptableObject
{
    public string ascendancyName = "New Ascendancy";
    // Assign the PassiveTreeData asset that defines the passive tree for this ascendancy class.
    public PassiveTreeData passiveTree; 

    // You can add other ascendancy-specific properties here, like unique mechanics, etc.
}
