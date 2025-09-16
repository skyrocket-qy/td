using UnityEngine;
using System.Collections.Generic;

[CreateAssetMenu(fileName = "FireballSkill", menuName = "Game/Skills/Fireball")]
public class FireballSkillData : SkillData
{
    public float burnDuration = 3f;
    public int burnDamagePerSecond = 1;
    public IgniteStatusEffectData igniteEffect; // Reference to an Ignite Status Effect asset

    public override void ExecuteSkill(Hero hero, List<Monster> targets)
    {
        Debug.Log($"{hero.name} casts Fireball!");
        base.ExecuteSkill(hero, targets); // Call base to apply initial damage

        foreach (Monster monster in targets)
        {
            if (igniteEffect != null)
            {
                monster.ApplyStatusEffect(igniteEffect);
            }
            else
            {
                Debug.LogWarning("Ignite Effect not assigned to Fireball Skill Data.");
            }
        }
    }
}
