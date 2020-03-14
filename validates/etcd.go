package validates

type CreateEtcdKVRequest struct {
	Key string `json:"key" validates:"required"  comment:"key1"`
	Value string `json:"value" validates:"required"  comment:"value"`
}
