package core_models

import "testing"

func TestToString(t *testing.T) {

	t.Run("Should return correct role name", func(t *testing.T) {
		actorRoles := map[int]string{
			0: "Lead Hero",
			1: "Lead Heroin",
			2: "Lead Billen",
			3: "Hero",
			4: "Heroin",
			5: "Billen",
			6: "Other",
		}
		for roleId, expectedRoleName := range actorRoles {
			model := ActorRole(roleId)

			roleName, err := model.ToString()

			if err != nil {
				t.Errorf("Expected '%s', but thrown error", expectedRoleName)
			}

			if roleName != expectedRoleName {
				t.Errorf("Got '%s', expected '%s'", roleName, expectedRoleName)
			}
		}
	})

	t.Run("Should return error for invalid role id", func(t *testing.T) {
		model := ActorRole(100)

		roleName, err := model.ToString()

		if err == nil {
			t.Error("Got '%s', expected error", roleName)
		}
	})
}
