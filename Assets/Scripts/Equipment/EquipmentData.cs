using UnityEngine;

namespace Equipment
{
    [CreateAssetMenu(fileName = "NewEquipment", menuName = "Game/Equipment Data")]
    public class EquipmentData : ScriptableObject
    {
        public string equipmentName = "New Equipment";
        [TextArea(3, 5)]
        public string description = "A brief description of the equipment's properties.";
        public Sprite icon; // Visual representation of the equipment

        public EquipmentSlotType slotType = EquipmentSlotType.Accessory; // Default to accessory

        // Stat bonuses this equipment provides
        public int bonusDamage = 0;
        public int bonusHealth = 0;
        public float bonusAttackRange = 0f;
        // Add more stats as needed (e.g., bonusDefense, bonusSpeed, etc.)
    }

    public enum EquipmentSlotType
    {
        Weapon,
        Armor,
        Accessory, // Can have multiple accessories
        // Add more specific slots if your game requires them (e.g., Helmet, Boots)
    }
}
