package RequestHandler



func (s APIFunctions) ServiceFinalize()(string,error){
	// テーブルをドロップする
	if err := s.db_repo.DropTable(); err != nil {
		return "Failed ServiceFinalize",err
	}
	// 今のところはドロップテーブルのみ
	return "Success ServiceFinalize",nil
}
