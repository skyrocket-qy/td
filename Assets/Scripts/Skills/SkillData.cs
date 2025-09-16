using UnityEngine;
using System.Collections.Generic;

[CreateAssetMenu(fileName = "NewSkillData", menuName = "Game/Skill Data")]
public class SkillData : ScriptableObject
{
    public string skillName = "New Skill";
    [TextArea] public string description = "Skill description here.";
    public Sprite icon; // Optional: for UI display

    public float cooldown = 5f;
    public float damageMultiplier = 1f; // Multiplier for hero's base damage
    public float rangeMultiplier = 1f;  // Multiplier for hero's base attack range

    public GameObject visualEffectPrefab; // Assign a particle system or other visual prefab here

    // Placeholder for skill-specific effects
    public virtual void ExecuteSkill(Hero hero, List<Monster> targets)
    {
        Debug.Log($"{hero.name} uses {skillName}!");

        // Instantiate visual effect if assigned
        if (visualEffectPrefab != null)
        {
            // For simplicity, instantiate at hero's position. You might want to instantiate at target or area center.
            GameObject effectInstance = Instantiate(visualEffectPrefab, hero.transform.position, Quaternion.identity);
            // Destroy the effect after a short duration (adjust as needed)
            Destroy(effectInstance, 2f); 
        }

        // Default skill logic (e.g., apply damage to targets)
        foreach (Monster monster in targets)
        {
            // Example: Apply damage based on hero's base damage and skill's damage multiplier
            int finalDamage = Mathf.RoundToInt(hero.damageAmount * damageMultiplier);
            monster.TakeDamage(finalDamage);
            Debug.Log($"{hero.name} dealt {finalDamage} skill damage to {monster.name}.");
        }
    }
}
