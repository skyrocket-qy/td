using UnityEngine;
using UnityEngine.AI; // Added for NavMeshAgent
using System.Collections.Generic; // Added for List
using System; // Added for Action delegate

public class Monster : MonoBehaviour
{
    public static event Action<Monster> OnMonsterDied; // Event to notify GameManager on death

    public int maxHealth = 10;
    public int currentHealth;

    private MonsterMovement monsterMovement;
    private StatusEffectManager statusEffectManager; // Reference to the new manager

    void Awake()
    {
        currentHealth = maxHealth;
        monsterMovement = GetComponent<MonsterMovement>();
        if (monsterMovement == null)
        {
            Debug.LogError("MonsterMovement component not found on this GameObject!");
        }

        statusEffectManager = GetComponent<StatusEffectManager>();
        if (statusEffectManager == null)
        {
            Debug.LogWarning("StatusEffectManager component not found on this Monster. Adding one.");
            statusEffectManager = gameObject.AddComponent<StatusEffectManager>();
        }
    }

    public void TakeDamage(int amount)
    {
        currentHealth -= amount;
        if (currentHealth <= 0)
        {
            Die();
        }
    }

    private void Die()
    {
        Debug.Log("Monster Died!");
        OnMonsterDied?.Invoke(this); // Notify GameManager that this monster has died
        // Add death animation, disable GameObject, etc.
        Destroy(gameObject);
    }

    // Example of how Monster might use MonsterMovement
    public void StartMoving()
    {
        if (monsterMovement != null)
        {
            // Assuming MonsterMovement has a method to start its movement
            // For now, let's just log that it would move.
            Debug.Log("Monster is commanded to start moving.");
            // monsterMovement.StartMovement(); // You would implement this in MonsterMovement
        }
    }

    // Public method to apply status effects to this monster
    public void ApplyStatusEffect(StatusEffectData effectData)
    {
        if (statusEffectManager != null)
        {
            statusEffectManager.ApplyStatusEffect(effectData);
        }
    }
}
