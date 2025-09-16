using UnityEngine;

[CreateAssetMenu(fileName = "NewStatusEffect", menuName = "Game/Status Effects/Base Status Effect")]
public class StatusEffectData : ScriptableObject
{
    public string effectName = "New Effect";
    [TextArea] public string description = "Description of the status effect.";
    public Sprite icon; // Optional: for UI display
    public float duration = 5f; // How long the effect lasts
    public float tickInterval = 0f; // How often the TickEffect is called (0 for non-ticking effects)

    // Virtual methods to be overridden by specific status effects
    public virtual void ApplyEffect(Monster target)
    {
        Debug.Log($"Applying {effectName} to {target.name}.");
        // Base logic for applying the effect (e.g., visual changes, initial stat modifications)
    }

    public virtual void TickEffect(Monster target)
    {
        // Logic for effects that happen over time (e.g., damage over time, healing over time)
    }

    public virtual void RemoveEffect(Monster target)
    {
        Debug.Log($"Removing {effectName} from {target.name}.");
        // Base logic for removing the effect (e.g., reverting visual changes, stat modifications)
    }
}
