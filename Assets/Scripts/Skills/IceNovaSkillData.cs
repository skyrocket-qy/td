using UnityEngine;
using System.Collections.Generic;

[CreateAssetMenu(fileName = "IceNovaSkill", menuName = "Game/Skills/Ice Nova")]
public class IceNovaSkillData : SkillData
{
    public float freezeDuration = 2f;
    public float slowAmount = 0.5f; // e.g., 50% slow
    public ChillStatusEffectData chillEffect; // Reference to a Chill Status Effect asset
    public FreezeStatusEffectData freezeEffect; // Reference to a Freeze Status Effect asset

    public override void ExecuteSkill(Hero hero, List<Monster> targets)
    {
        Debug.Log($"{hero.name} casts Ice Nova!");
        base.ExecuteSkill(hero, targets); // Call base to apply initial damage

        foreach (Monster monster in targets)
        {
            if (chillEffect != null)
            {
                monster.ApplyStatusEffect(chillEffect);
            }
            else
            {
                Debug.LogWarning("Chill Effect not assigned to Ice Nova Skill Data.");
            }

            // Optionally, apply freeze if conditions are met (e.g., monster is already chilled)
            if (freezeEffect != null)
            {
                monster.ApplyStatusEffect(freezeEffect);
            }
            else
            {
                Debug.LogWarning("Freeze Effect not assigned to Ice Nova Skill Data.");
            }
        }
    }
}
