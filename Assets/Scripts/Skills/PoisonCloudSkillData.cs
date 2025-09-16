using UnityEngine;
using System.Collections.Generic;

[CreateAssetMenu(fileName = "PoisonCloudSkill", menuName = "Game/Skills/Poison Cloud")]
public class PoisonCloudSkillData : SkillData
{
    public float poisonDuration = 4f;
    public int poisonDamagePerTick = 1;
    public float tickInterval = 1f;
    public PoisonStatusEffectData poisonEffect; // Reference to a Poison Status Effect asset

    public override void ExecuteSkill(Hero hero, List<Monster> targets)
    {
        Debug.Log($"{hero.name} casts Poison Cloud!");
        base.ExecuteSkill(hero, targets); // Call base to apply initial damage

        foreach (Monster monster in targets)
        {
            if (poisonEffect != null)
            {
                monster.ApplyStatusEffect(poisonEffect);
            }
            else
            {
                Debug.LogWarning("Poison Effect not assigned to Poison Cloud Skill Data.");
            }
        }
    }
}
