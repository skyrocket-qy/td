using UnityEngine;
using System.Collections.Generic;
using System.Linq; // Added for LINQ operations
using Equipment; // Added to access EquipmentData and EquipmentSlotType

public class Hero : MonoBehaviour
{
    // Base Stats (not modified directly by equipment, only by level/passives)
    public int baseDamageAmount = 5;
    public float baseAttackRange = 2f;

    // Current Stats (derived from base stats + equipment + passives)
    public int damageAmount;
    public float attackRange;

    // Hero Progression Fields
    public int level = 1;
    public HeroClassData heroClassData; // Now references a ScriptableObject
    public AscendancyClassData ascendancyClassData; // Now references a ScriptableObject

    // Experience and Leveling
    public int currentExperience = 0;
    public int experienceToNextLevel = 100; // Example value, adjust as needed

    // Equipment Fields
    public EquipmentData[] equippedItems = new EquipmentData[3]; // A hero can equip 3 items

    public void GainExperience(int amount)
    {
        currentExperience += amount;
        Debug.Log($"Gained {amount} experience. Current experience: {currentExperience}");

        while (currentExperience >= experienceToNextLevel && level < 20)
        {
            LevelUp();
        }
    }

    private void LevelUp()
    {
        level++;
        currentExperience -= experienceToNextLevel;
        experienceToNextLevel = CalculateExperienceToNextLevel(level); // Implement this method
        Debug.Log($"Hero leveled up to level {level}!");

        if (level == 10 && ascendancyClassData == null)
        {
            Debug.Log("Congratulations! You have reached level 10. You can now choose an Ascendancy Class.");
            // In a real game, this would trigger a UI to select the Ascendancy Class.
            // For now, we'll just log a message.
        }

        RecalculateStats(); // Recalculate stats after leveling up
    }

    private int CalculateExperienceToNextLevel(int currentLevel)
    {
        // Simple example: experience needed increases with level
        return 100 + (currentLevel * 50);
    }

    // Method to assign Ascendancy Class (would be called from UI after selection)
    public void AssignAscendancyClass(AscendancyClassData newAscendancyClass)
    {
        if (level >= 10 && ascendancyClassData == null)
        {
            ascendancyClassData = newAscendancyClass;
            Debug.Log($"Ascendancy Class '{newAscendancyClass.ascendancyName}' assigned!");
            RecalculateStats(); // Recalculate stats after assigning ascendancy
        }
        else if (ascendancyClassData != null)
        {
            Debug.LogWarning("Ascendancy Class already assigned.");
        }
        else
        {
            Debug.LogWarning("Cannot assign Ascendancy Class before reaching level 10.");
        }
    }

    // Equipment Management
    public bool Equip(EquipmentData newEquipment)
    {
        if (newEquipment == null)
        {
            Debug.LogWarning("Attempted to equip null equipment.");
            return false;
        }

        // Try to find an empty slot
        for (int i = 0; i < equippedItems.Length; i++)
        {
            if (equippedItems[i] == null)
            {
                equippedItems[i] = newEquipment;
                Debug.Log($"Equipped {newEquipment.equipmentName} in slot {i}.");
                RecalculateStats();
                return true;
            }
        }

        // If no empty slot, try to replace an item of the same slot type
        for (int i = 0; i < equippedItems.Length; i++)
        {
            if (equippedItems[i] != null && equippedItems[i].slotType == newEquipment.slotType)
            {
                Debug.Log($"Replacing {equippedItems[i].equipmentName} with {newEquipment.equipmentName} in slot {i}.");
                equippedItems[i] = newEquipment;
                RecalculateStats();
                return true;
            }
        }

        Debug.LogWarning($"Could not equip {newEquipment.equipmentName}. All slots are full and no matching slot type found for replacement.");
        return false;
    }

    public bool Unequip(EquipmentData equipmentToRemove)
    {
        if (equipmentToRemove == null)
        {
            Debug.LogWarning("Attempted to unequip null equipment.");
            return false;
        }

        for (int i = 0; i < equippedItems.Length; i++)
        {
            if (equippedItems[i] == equipmentToRemove)
            {
                equippedItems[i] = null;
                Debug.Log($"Unequipped {equipmentToRemove.equipmentName} from slot {i}.");
                RecalculateStats();
                return true;
            }
        }
        Debug.LogWarning($"Could not unequip {equipmentToRemove.equipmentName}. Item not found in equipped slots.");
        return false;
    }

    // Skill Mechanism Fields
    public List<SkillData> heroSkills = new List<SkillData>();
    private Dictionary<SkillData, float> skillCooldowns = new Dictionary<SkillData, float>();

    void Awake()
    {
        foreach (SkillData skill in heroSkills)
        {
            skillCooldowns.Add(skill, 0f);
        }

        RecalculateStats(); // Initial stat calculation on Awake
    }

    void Update()
    {
        // Update skill cooldowns
        List<SkillData> skillsOnCooldown = new List<SkillData>();
        foreach (var entry in skillCooldowns)
        {
            if (entry.Value > 0)
            {
                skillsOnCooldown.Add(entry.Key);
            }
        }

        foreach (SkillData skill in skillsOnCooldown)
        {
            skillCooldowns[skill] -= Time.deltaTime;
        }

        // Example: Press 'Space' to attack
        if (Input.GetKeyDown(KeyCode.Space))
        {
            Attack();
        }

        // Example: Press '1' to use the first skill (if available)
        if (Input.GetKeyDown(KeyCode.Alpha1) && heroSkills.Count > 0)
        {
            UseSkill(heroSkills[0]);
        }
    }

    void Attack()
    {
        Debug.Log("Hero attacks!");

        // Find all monsters in the scene (for simplicity)
        // In a real game, you'd likely use physics overlaps or raycasts
        // to find monsters within attack range.
        Monster[] monsters = FindObjectsByType<Monster>(FindObjectsSortMode.None);

        foreach (Monster monster in monsters)
        {
            // Check if the monster is within attack range
            if (Vector3.Distance(transform.position, monster.transform.position) <= attackRange)
            {
                monster.TakeDamage(damageAmount);
                Debug.Log($"Hero dealt {damageAmount} damage to {monster.name}.");
                // For demonstration, we'll only hit one monster per attack for now.
                // You might want to hit multiple in an AoE attack.
                return; 
            }
        }
        Debug.Log("No monster found within attack range.");
    }

    // Method for using a skill
    public void UseSkill(SkillData skill)
    {
        if (skill == null) return;

        if (skillCooldowns.ContainsKey(skill) && skillCooldowns[skill] <= 0)
        {
            Debug.Log($"Hero attempts to use {skill.skillName}!");

            // Find targets within the skill's effective range
            List<Monster> targets = new List<Monster>();
            float effectiveSkillRange = attackRange * skill.rangeMultiplier;
            Monster[] allMonsters = FindObjectsByType<Monster>(FindObjectsSortMode.None);

            foreach (Monster monster in allMonsters)
            {
                if (Vector3.Distance(transform.position, monster.transform.position) <= effectiveSkillRange)
                {
                    targets.Add(monster);
                }
            }

            if (targets.Count > 0)
            {
                skill.ExecuteSkill(this, targets);
                skillCooldowns[skill] = skill.cooldown; // Reset cooldown
            }
            else
            {
                Debug.Log($"No targets found for {skill.skillName} within range.");
            }
        }
        else if (skillCooldowns.ContainsKey(skill))
        {
            Debug.Log($"{skill.skillName} is on cooldown. Remaining: {skillCooldowns[skill]:F2}s");
        }
        else
        {
            Debug.LogWarning($"Skill {skill.skillName} not found in hero's skill list.");
        }
    }

    // Placeholder method to apply passive abilities from assigned classes
    // TODO: Implement actual stat application here.
    private void ApplyPassiveAbilities()
    {
        // This method should ideally just trigger a stat recalculation, not apply stats directly.
        // The actual application of passive bonuses is handled in ApplyPassiveBonuses() called by RecalculateStats().
        // Any direct stat manipulation here would be reset by RecalculateStats().
        RecalculateStats(); 
    }

    private void RecalculateStats()
    {
        // Reset to base stats
        damageAmount = baseDamageAmount;
        attackRange = baseAttackRange;

        // Apply bonuses from equipped items
        foreach (EquipmentData equipment in equippedItems)
        {
            if (equipment != null)
            {
                damageAmount += equipment.bonusDamage;
                attackRange += equipment.bonusAttackRange;
                // Apply other equipment bonuses here (e.g., health, defense)
            }
        }

        // Apply bonuses from passive abilities
        ApplyPassiveBonuses();

        Debug.Log($"Hero Stats Recalculated: Damage={damageAmount}, Range={attackRange}");
    }

    // Helper method to apply only passive bonuses, called by RecalculateStats
    private void ApplyPassiveBonuses()
    {
        if (heroClassData != null && heroClassData.passiveTree != null)
        {
            foreach (PassiveAbilityData ability in heroClassData.passiveTree.unlockedAbilities)
            {
                damageAmount += ability.bonusDamage;
                attackRange += ability.bonusAttackRange;
            }
        }

        if (ascendancyClassData != null && ascendancyClassData.passiveTree != null)
        {
            foreach (PassiveAbilityData ability in ascendancyClassData.passiveTree.unlockedAbilities)
            {
                damageAmount += ability.bonusDamage;
                attackRange += ability.bonusAttackRange;
            }
        }
    }
}